package gui

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	ErrEmptyJSON = errors.New("empty json")
)

type Gui struct {
	Tree  *Tree
	Navi  *Navi
	App   *tview.Application
	Pages *tview.Pages
}

func New() *Gui {
	g := &Gui{
		Tree:  NewTree(),
		Navi:  NewNavi(),
		App:   tview.NewApplication(),
		Pages: tview.NewPages(),
	}
	return g
}

func (g *Gui) Run(i interface{}) error {
	g.Tree.UpdateView(g, i)
	g.Tree.SetKeybindings(g)
	g.Navi.UpdateView()
	g.Navi.SetKeybindings(g)

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
			g.Pages.SwitchToPage(page).ShowPage("main")
			if buttonLabel == doneLabel {
				doneFunc()
			}
		})

	g.Pages.AddAndSwitchToPage("message", g.Modal(modal, 80, 29), true).ShowPage("main")
}

func (g *Gui) Input(text, label string, width int, doneFunc func(text string)) {
	input := tview.NewInputField().SetText(text)
	input.SetBorder(true)
	input.SetLabel(label).SetLabelWidth(width).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			doneFunc(input.GetText())
			g.Pages.RemovePage("input")
		}
	})

	g.Pages.AddAndSwitchToPage("input", g.Modal(input, 0, 3), true).ShowPage("main")
}

func (g *Gui) Form(fieldLabel []string, doneLabel, title, pageName string,
	height int, doneFunc func(values map[string]string) error) {

	if g.Pages.HasPage(pageName) {
		g.Pages.ShowPage(pageName)
		return
	}

	form := tview.NewForm()
	for _, label := range fieldLabel {
		form.AddInputField(label, "", 0, nil, nil)
	}

	form.AddButton(doneLabel, func() {
		values := make(map[string]string)

		for _, label := range fieldLabel {
			item := form.GetFormItemByLabel(label)
			switch item.(type) {
			case *tview.InputField:
				input, ok := item.(*tview.InputField)
				if ok {
					values[label] = os.ExpandEnv(input.GetText())
				}
			}
		}

		if err := doneFunc(values); err != nil {
			g.Message(err.Error(), pageName, func() {})
			return
		}

		g.Pages.RemovePage(pageName)
	}).
		AddButton("cancel", func() {
			g.Pages.RemovePage(pageName)
		})

	form.SetBorder(true).SetTitle(title).
		SetTitleAlign(tview.AlignLeft)

	g.Pages.AddAndSwitchToPage(pageName, g.Modal(form, 0, height), true).ShowPage("main")
}

func (g *Gui) LoadJSON() {
	labels := []string{"file"}
	g.Form(labels, "read", "read from file", "read_from_file", 7, func(values map[string]string) error {
		fileName := values[labels[0]]
		file, err := os.Open(fileName)
		if err != nil {
			log.Println(fmt.Sprintf("can't open file: %s", err))
			return err
		}

		i, err := UnMarshalJSON(file)
		if err != nil {
			return err
		}

		g.Tree.UpdateView(g, i)
		return nil
	})
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
				root.SetChildren(g.walk(root.GetChildren(), text))
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

func (g *Gui) walk(nodes []*tview.TreeNode, text string) []*tview.TreeNode {
	var newNodes []*tview.TreeNode

	for _, child := range nodes {
		log.Println(child.GetText())
		if strings.Index(strings.ToLower(child.GetText()), strings.ToLower(text)) != -1 {
			newNodes = append(newNodes, child)
		} else {
			newNodes = append(newNodes, g.walk(child.GetChildren(), text)...)
		}
	}

	return newNodes
}

func (g *Gui) SaveJSON() {
	labels := []string{"file"}
	g.Form(labels, "save", "save to file", "save_to_file", 7, func(values map[string]string) error {
		file := values[labels[0]]
		file = os.ExpandEnv(file)
		return g.SaveJSONToFile(file)
	})
}

func (g *Gui) SaveJSONToFile(file string) error {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")

	if err := enc.Encode(g.makeJSON(g.Tree.GetRoot())); err != nil {
		log.Println(fmt.Sprintf("can't marshal json: %s", err))
		return err
	}

	if err := ioutil.WriteFile(file, buf.Bytes(), 0666); err != nil {
		log.Println(fmt.Sprintf("can't create file: %s", err))
		return err
	}

	return nil
}

func (g *Gui) makeJSON(node *tview.TreeNode) interface{} {
	ref := node.GetReference().(Reference)
	children := node.GetChildren()

	switch ref.JSONType {
	case Object:
		i := make(map[string]interface{})
		for _, n := range children {
			i[n.GetText()] = g.makeJSON(n)
		}
		return i
	case Array:
		var i []interface{}
		for _, n := range children {
			i = append(i, g.makeJSON(n))
		}
		return i
	case Key:
		v := node.GetChildren()[0]
		if v.GetReference().(Reference).JSONType == Value {
			return g.parseValue(v)
		}
		return map[string]interface{}{
			node.GetText(): g.makeJSON(v),
		}
	}

	return g.parseValue(node)
}

func (g *Gui) parseValue(node *tview.TreeNode) interface{} {
	v := node.GetText()
	ref := node.GetReference().(Reference)

	switch ref.ValueType {
	case Int:
		i, _ := strconv.Atoi(v)
		return i
	case Float:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	case Boolean:
		b, _ := strconv.ParseBool(v)
		return b
	case Null:
		return nil
	}

	return v
}

func (g *Gui) AddNode() {
	labels := []string{"json"}
	g.Form(labels, "add", "add new node", "add_new_node", 7, func(values map[string]string) error {
		j := values[labels[0]]
		if j == "" {
			log.Println(ErrEmptyJSON)
			return ErrEmptyJSON
		}

		buf := bytes.NewBufferString(j)
		i, err := UnMarshalJSON(buf)
		if err != nil {
			return err
		}

		newNode := NewRootTreeNode(i)
		newNode.SetChildren(g.Tree.AddNode(i))
		g.Tree.GetCurrentNode().AddChild(newNode)
		// update new origin root node
		g.Tree.OriginRoot = g.Tree.GetRoot()

		return nil
	})
}

func (g *Gui) AddValue() {
	labels := []string{"json"}
	g.Form(labels, "add", "add new value", "add_new_value", 7, func(values map[string]string) error {
		j := values[labels[0]]
		if j == "" {
			log.Println(ErrEmptyJSON)
			return ErrEmptyJSON
		}

		buf := bytes.NewBufferString(j)
		i, err := UnMarshalJSON(buf)
		if err != nil {
			return err
		}

		current := g.Tree.GetCurrentNode()
		for _, n := range g.Tree.AddNode(i) {
			current.AddChild(n)
		}
		// update new origin root node
		g.Tree.OriginRoot = g.Tree.GetRoot()

		return nil
	})
}

func (g *Gui) NaviPanel() {
	if g.Pages.HasPage(NaviPageName) {
		g.Pages.ShowPage(NaviPageName)
	} else {
		g.Pages.AddAndSwitchToPage(NaviPageName, g.Modal(g.Navi, 0, 0), true).ShowPage("main")
	}
}

func (g *Gui) EditWithEditor() {
	f, err := ioutil.TempFile("", "tson")
	if err != nil {
		log.Println(fmt.Sprintf("can't create temp file: %s", err))
		g.Message(err.Error(), "main", func() {})
		return
	}

	defer os.RemoveAll(f.Name())

	if err := g.SaveJSONToFile(f.Name()); err != nil {
		log.Println(fmt.Sprintf("can't write to temp file: %s", err))
		g.Message(err.Error(), "main", func() {})
		return
	}

	if b := g.App.Suspend(func() {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			log.Println(fmt.Sprintf("$EDITOR is empty: %s", err))
			g.Message(err.Error(), "main", func() {})
			return
		}

		cmd := exec.Command(editor, f.Name())

		// Start the command with a pty.
		ptmx, err := pty.Start(cmd)
		if err != nil {
			g.Message(err.Error(), "main", func() {})
			return
		}
		// Make sure to close the pty at the end.
		defer func() { _ = ptmx.Close() }() // Best effort.

		// Handle pty size.
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGWINCH)
		go func() {
			for range ch {
				if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
					log.Printf("can't resizing pty: %s", err)
				}
			}
		}()
		ch <- syscall.SIGWINCH // Initial resize.

		// Set stdin in raw mode.
		oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			log.Println(fmt.Sprintf("can't make terminal raw mode: %s", err))
			g.Message(err.Error(), "main", func() {})
			return
		}
		defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

		// Copy stdin to the pty and the pty to stdout.
		go func() {
			io.Copy(ptmx, os.Stdin)
		}()

		io.Copy(os.Stdout, ptmx)

		i, err := UnMarshalJSON(f)
		if err != nil {
			log.Println(fmt.Sprintf("can't read from file: %s", err))
			g.Message(err.Error(), "main", func() {})
			return
		}

		g.Tree.UpdateView(g, i)
	}); !b {
		log.Println(fmt.Sprintf("can't edit: %s", err))
		g.Message(err.Error(), "main", func() {})
		return
	}
}

func UnMarshalJSON(in io.Reader) (interface{}, error) {
	b, err := ioutil.ReadAll(in)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(b) == 0 {
		log.Println(err)
		return nil, ErrEmptyJSON
	}

	var i interface{}
	if err := json.Unmarshal(b, &i); err != nil {
		log.Println(err)
		return nil, err
	}

	return i, nil
}
