package gui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Gui struct {
	Tree  *Tree
	App   *tview.Application
	Pages *tview.Pages
}

func New() *Gui {
	g := &Gui{
		Tree:  NewTree(),
		App:   tview.NewApplication(),
		Pages: tview.NewPages(),
	}
	return g
}

func (g *Gui) Run(i interface{}) error {
	g.Tree.UpdateView(g, i)
	g.Tree.SetKeybindings(g)

	grid := tview.NewGrid().
		AddItem(g.Tree, 0, 0, 1, 1, 0, 0, true)

	g.Pages.AddAndSwitchToPage("main", grid, true)

	if err := g.App.SetRoot(g.Pages, true).Run(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (g *Gui) Modal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewGrid().
		SetColumns(0, width, 0).
		SetRows(0, height, 0).
		AddItem(p, 1, 1, 1, 1, 0, 0, true)
}

func (g *Gui) Message(message, page string, doneFunc func()) {
	doneLabel := "ok"
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{doneLabel}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			g.Pages.RemovePage("message")
			g.Pages.SwitchToPage(page)
			if buttonLabel == doneLabel {
				doneFunc()
			}
		})

	g.Pages.AddAndSwitchToPage("message", g.Modal(modal, 80, 29), true).ShowPage("main")
}

func (g *Gui) Input(text, label string, doneFunc func(text string)) {
	input := tview.NewInputField().SetText(text)
	input.SetBorder(true)
	input.SetLabel(label).SetLabelWidth(7).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			doneFunc(input.GetText())
			g.Pages.RemovePage("input")
		}
	})

	g.Pages.AddAndSwitchToPage("input", g.Modal(input, 0, 3), true).ShowPage("main")
}

func (g *Gui) LoadJSON() {
	pageName := "read_from_file"
	form := tview.NewForm()
	form.AddInputField("file", "", 0, nil, nil).
		AddButton("read", func() {
			file := form.GetFormItem(0).(*tview.InputField).GetText()
			file = os.ExpandEnv(file)
			b, err := ioutil.ReadFile(file)
			if err != nil {
				msg := fmt.Sprintf("can't read file: %s", err)
				log.Println(msg)
				g.Message(msg, "main", func() {})
				return
			}

			var i interface{}
			if err := json.Unmarshal(b, &i); err != nil {
				msg := fmt.Sprintf("can't read file: %s", err)
				log.Println(msg)
				g.Message(msg, "main", func() {})
				return
			}

			g.Tree.UpdateView(g, i)
			g.Pages.RemovePage(pageName)
		}).
		AddButton("cancel", func() {
			g.Pages.RemovePage(pageName)
		})

	form.SetBorder(true).SetTitle("read from file").
		SetTitleAlign(tview.AlignLeft)

	g.Pages.AddAndSwitchToPage(pageName, g.Modal(form, 0, 8), true).ShowPage("main")
}

func (g *Gui) Search() {
	pageName := "search"
	if g.Pages.HasPage(pageName) {
		g.Pages.ShowPage(pageName)
	} else {
		input := tview.NewInputField()
		input.SetBorder(true).SetTitle("search").SetTitleAlign(tview.AlignLeft)
		input.SetChangedFunc(func(text string) {
			root := *g.Tree.OriginRoot
			g.Tree.SetRoot(&root)
			if text != "" {
				root := g.Tree.GetRoot()
				root.SetChildren(g.walk(root, text))
			}
		})
		input.SetLabel("word").SetLabelWidth(5).SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				g.Pages.HidePage(pageName)
			}
		})

		g.Pages.AddAndSwitchToPage(pageName, g.Modal(input, 0, 3), true).ShowPage("main")
	}
}

func (g *Gui) walk(node *tview.TreeNode, text string) []*tview.TreeNode {
	var nodes []*tview.TreeNode
	if strings.Index(strings.ToLower(node.GetText()), text) != -1 {
		nodes = append(nodes, node)
		return nodes
	}

	for _, node := range node.GetChildren() {
		nodes = append(nodes, g.walk(node, text)...)
	}

	return nodes

}

func (g *Gui) SaveJSON() {
	pageName := "save_to_file"
	form := tview.NewForm()
	form.AddInputField("file", "", 0, nil, nil).
		AddButton("save", func() {
			fileName := form.GetFormItem(0).(*tview.InputField).GetText()
			fileName = os.ExpandEnv(fileName)

			var b bytes.Buffer
			enc := json.NewEncoder(&b)
			enc.SetIndent("", "  ")

			if err := enc.Encode(g.Tree.OriginJSON); err != nil {
				msg := fmt.Sprintf("can't make json: %s", err)
				log.Println(msg)
				g.Message(msg, "main", func() {})
				return
			}

			if err := ioutil.WriteFile(fileName, b.Bytes(), 0666); err != nil {
				msg := fmt.Sprintf("can't create file: %s", err)
				log.Println(msg)
				g.Message(msg, "main", func() {})
				return
			}
			g.Pages.RemovePage(pageName)
		}).
		AddButton("cancel", func() {
			g.Pages.RemovePage(pageName)
		})

	form.SetBorder(true).SetTitle("save to file").
		SetTitleAlign(tview.AlignLeft)

	g.Pages.AddAndSwitchToPage(pageName, g.Modal(form, 0, 8), true).ShowPage("main")
}
