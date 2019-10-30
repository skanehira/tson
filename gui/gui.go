package gui

import (
	"log"

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

func (g *Gui) Input(text string, doneFunc func(text string)) {
	input := tview.NewInputField().SetText(text)
	input.SetLabel("field:").SetLabelWidth(6).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			doneFunc(input.GetText())
			g.Pages.SendToBack("input")
		}
	})

	g.Pages.AddPage("input", input, true, true).SendToFront("input").ShowPage("main")
}
