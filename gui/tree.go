package gui

import (
	"fmt"
	"log"
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
		r := reflect.ValueOf(i)

		var root *tview.TreeNode
		switch r.Kind() {
		case reflect.Map:
			root = tview.NewTreeNode("{object}").SetReference(Reference{JSONType: Object})
		case reflect.Slice:
			root = tview.NewTreeNode("{array}").SetReference(Reference{JSONType: Array})
		default:
			root = tview.NewTreeNode("{value}").SetReference(Reference{JSONType: Key})
		}

		root.SetChildren(t.AddNode(i))
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
				SetColor(tcell.ColorMediumSlateBlue).
				SetChildren(t.AddNode(v))
			r := reflect.ValueOf(v)

			if r.Kind() == reflect.Slice {
				newNode.SetReference(Reference{JSONType: Array})
			} else if r.Kind() == reflect.Map {
				newNode.SetReference(Reference{JSONType: Object})
			} else {
				newNode.SetReference(Reference{JSONType: Key})
			}

			log.Printf("key:%v value:%v value_kind:%v", k, v, newNode.GetReference())
			nodes = append(nodes, newNode)
		}
	case []interface{}:
		for _, v := range node {
			switch n := v.(type) {
			case map[string]interface{}:
				r := reflect.ValueOf(n)
				if r.Kind() != reflect.Slice {
					objectNode := tview.NewTreeNode("{object}").
						SetChildren(t.AddNode(v)).SetReference(Reference{JSONType: Object})

					log.Printf("value:%v value_kind:%v", v, "object")
					nodes = append(nodes, objectNode)
				}
			default:
				nodes = append(nodes, t.AddNode(v)...)
			}
		}
	default:
		log.Printf("value:%v value_kind:%v", node, "value")
		ref := reflect.ValueOf(node)
		var valueType ValueType
		switch ref.Kind() {
		case reflect.Int:
			valueType = Int
		case reflect.Float64:
			valueType = Float
		case reflect.Bool:
			valueType = Boolean
		default:
			if node == nil {
				valueType = Null
			} else {
				valueType = String
			}
		}

		log.Printf("value_type:%v", valueType)
		nodes = append(nodes, t.NewNodeWithLiteral(node).
			SetReference(Reference{JSONType: Value, ValueType: valueType}))
	}
	return nodes
}

func (t *Tree) NewNodeWithLiteral(i interface{}) *tview.TreeNode {
	if i == nil {
		return tview.NewTreeNode("null")
	}
	return tview.NewTreeNode(fmt.Sprintf("%v", i))
}

func (t *Tree) SetKeybindings(g *Gui) {
	t.SetSelectedFunc(func(node *tview.TreeNode) {
		text := node.GetText()
		if text == "{object}" || text == "{array}" || text == "{value}" {
			return
		}
		labelWidth := 5
		g.Input(text, "text", labelWidth, func(text string) {
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
		case 's':
			g.SaveJSON()
		case '/':
			g.Search()
		}

		return event
	})
}
