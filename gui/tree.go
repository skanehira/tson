package gui

import (
	"fmt"
	"reflect"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Tree struct {
	*tview.TreeView
	OriginRoot *tview.TreeNode
}

func NewTree() *Tree {
	t := &Tree{
		TreeView: tview.NewTreeView(),
	}

	t.SetBorder(true).SetTitle("json tree").SetTitleAlign(tview.AlignLeft)
	return t
}

func (t *Tree) UpdateView(g *Gui, i interface{}) {
	g.App.QueueUpdateDraw(func() {
		root := tview.NewTreeNode(".").SetChildren(t.AddNode(i))
		t.SetRoot(root).SetCurrentNode(root)
		originRoot := *root
		t.OriginRoot = &originRoot
	})
}

func (t *Tree) AddNode(node interface{}) []*tview.TreeNode {
	var nodes []*tview.TreeNode

	switch node := node.(type) {
	case map[string]interface{}:
		for k, v := range node {
			newNode := t.NewNodeWithLiteral(k).
				SetColor(tcell.ColorMediumSlateBlue).SetReference(k)

			list, isList := v.([]interface{})
			if isList && len(list) > 0 {
				newNode.SetSelectable(true)
			}
			newNode.SetChildren(t.AddNode(v))
			nodes = append(nodes, newNode)
		}
	case []interface{}:
		for i, v := range node {
			switch n := v.(type) {
			case map[string]interface{}, []interface{}:
				if reflect.ValueOf(n).Len() > 0 {
					numberNode := tview.NewTreeNode(fmt.Sprintf("[%d]", i+1))
					numberNode.SetChildren(t.AddNode(v))
					nodes = append(nodes, numberNode)
				}
			default:
				nodes = append(nodes, t.AddNode(v)...)
			}
		}
	default:
		nodes = append(nodes, t.NewNodeWithLiteral(node))
	}
	return nodes
}

func (t *Tree) NewNodeWithLiteral(i interface{}) *tview.TreeNode {
	return tview.NewTreeNode(fmt.Sprintf("%v", i))
}

func (t *Tree) SetKeybindings(g *Gui) {
	t.SetSelectedFunc(func(node *tview.TreeNode) {
		g.Input(node.GetText(), "filed", func(text string) {
			node.SetText(text)
		})
	})

	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'h':
			t.GetCurrentNode().SetExpanded(false)
		case 'H':
			t.GetRoot().CollapseAll()
		case 'd':
			t.GetCurrentNode().ClearChildren()
		case 'L':
			t.GetRoot().ExpandAll()
		case 'l':
			t.GetCurrentNode().SetExpanded(true)
		case 'r':
			g.LoadJSON()
		case '/':
			g.Search()
		}

		return event
	})
}
