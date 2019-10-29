package gui

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Tree struct {
	*tview.TreeView
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
		root := tview.NewTreeNode(".")
		t.SetRoot(root).SetCurrentNode(root)
		for _, node := range t.AddNode(i) {
			root.AddChild(node)
		}
	})
}

func (t *Tree) AddNode(node interface{}) []*tview.TreeNode {
	// e.g child is {"name": "gorilla", "lang": {"ja":"japan", "en": "english"}}
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
			for _, n := range t.AddNode(v) {
				newNode.AddChild(n)
			}
			nodes = append(nodes, newNode)
		}
	case []interface{}:
		for i, v := range node {
			if list, isList := v.([]interface{}); isList && len(list) > 0 {
				numberNode := tview.NewTreeNode(fmt.Sprintf("[%d]", i+1))
				for _, n := range t.AddNode(v) {
					numberNode.AddChild(n)
				}
				nodes = append(nodes, numberNode)
			} else if m, isMap := v.(map[string]interface{}); isMap && len(m) > 0 {
				numberNode := tview.NewTreeNode(fmt.Sprintf("[%d]", i+1))
				for _, n := range t.AddNode(v) {
					numberNode.AddChild(n)
				}
				nodes = append(nodes, numberNode)
			} else {
				nodes = append(nodes, t.AddNode(v)...)
			}
		}
	default:
		nodes = append(nodes, t.NewNodeWithLiteral(node))
	}
	return nodes
}

func (t *Tree) NewNodeWithLiteral(i interface{}) *tview.TreeNode {
	var text string
	node := tview.NewTreeNode("")
	switch v := i.(type) {
	case int32:
		text = fmt.Sprintf("%d", v)
	case int64:
		text = fmt.Sprintf("%d", v)
	case float32:
		text = fmt.Sprintf("%f", v)
	case float64:
		text = fmt.Sprintf("%f", v)
	case bool:
		text = fmt.Sprintf("%t", v)
	case string:
		text = v
	}

	return node.SetText(text)
}
