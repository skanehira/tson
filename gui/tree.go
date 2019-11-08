package gui

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/gofrs/uuid"
	"github.com/rivo/tview"
)

const (
	moveToNext int = iota + 1
	moveToPre
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

		root := NewRootTreeNode(i)
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

			id := uuid.Must(uuid.NewV4()).String()
			if r.Kind() == reflect.Slice {
				newNode.SetReference(Reference{ID: id, JSONType: Array})
			} else if r.Kind() == reflect.Map {
				newNode.SetReference(Reference{ID: id, JSONType: Object})
			} else {
				newNode.SetReference(Reference{ID: id, JSONType: Key})
			}

			nodes = append(nodes, newNode)
		}
	case []interface{}:
		for _, v := range node {
			id := uuid.Must(uuid.NewV4()).String()
			switch v.(type) {
			case map[string]interface{}:
				objectNode := tview.NewTreeNode("{object}").
					SetChildren(t.AddNode(v)).SetReference(Reference{ID: id, JSONType: Object})
				nodes = append(nodes, objectNode)
			case []interface{}:
				arrayNode := tview.NewTreeNode("{array}").
					SetChildren(t.AddNode(v)).SetReference(Reference{ID: id, JSONType: Array})
				nodes = append(nodes, arrayNode)
			default:
				nodes = append(nodes, t.AddNode(v)...)
			}
		}
	default:
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

		id := uuid.Must(uuid.NewV4()).String()
		nodes = append(nodes, t.NewNodeWithLiteral(node).
			SetReference(Reference{ID: id, JSONType: Value, ValueType: valueType}))
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
			ref := node.GetReference().(Reference)
			ref.ValueType = parseValueType(text)
			if ref.ValueType == String {
				text = strings.Trim(text, `"`)
			}
			node.SetText(text).SetReference(ref)
		})
	})

	t.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'h':
			t.GetCurrentNode().SetExpanded(false)
		case 'H':
			t.CollapseValues(t.GetRoot())
		case 'd':
			t.GetCurrentNode().ClearChildren()
			newRoot := *g.Tree.GetRoot()
			g.Tree.OriginRoot = &newRoot
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
		case 'a':
			g.AddNode()
		case 'A':
			g.AddValue()
		case '?':
			g.NaviPanel()
		case 'e':
			g.EditWithEditor()
		case ' ':
			current := t.GetCurrentNode()
			current.SetExpanded(!current.IsExpanded())
		case 'q':
			g.App.Stop()
		}

		switch event.Key() {
		case tcell.KeyCtrlJ:
			t.moveParent(moveToNext)
		case tcell.KeyCtrlK:
			t.moveParent(moveToPre)
		}

		return event
	})
}

func (t *Tree) CollapseValues(node *tview.TreeNode) {
	node.Walk(func(node, parent *tview.TreeNode) bool {
		ref := node.GetReference().(Reference)
		if ref.JSONType == Value {
			pRef := parent.GetReference().(Reference)
			t := pRef.JSONType
			if t == Key || t == Array {
				parent.SetExpanded(false)
			}
		}
		return true
	})
}

func (t *Tree) moveParent(movement int) {
	current := t.GetCurrentNode()
	t.GetRoot().Walk(func(node, parent *tview.TreeNode) bool {
		if parent != nil {
			children := parent.GetChildren()
			for i, n := range children {
				if n.GetReference().(Reference).ID == current.GetReference().(Reference).ID {
					if movement == moveToNext {
						if i < len(children)-1 {
							t.SetCurrentNode(children[i+1])
							return false
						}
					} else if movement == moveToPre {
						if i > 0 {
							t.SetCurrentNode(children[i-1])
							return false
						}
					}
				}
			}
		}

		return true
	})
}

func parseValueType(text string) ValueType {
	// if sorround with `"` set string type
	if strings.HasPrefix(text, `"`) && strings.HasSuffix(text, `"`) {
		return String
	} else if "null" == text {
		return Null
	} else if text == "false" || text == "true" {
		return Boolean
	} else if _, err := strconv.ParseFloat(text, 64); err == nil {
		return Float
	} else if _, err := strconv.Atoi(text); err == nil {
		return Int
	}

	return String
}

func NewRootTreeNode(i interface{}) *tview.TreeNode {
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
	return root
}
