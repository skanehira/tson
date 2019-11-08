package gui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var NaviPageName = "navi_panel"

var RedColor = `[red::b]%s[white]: %s`

// default keybinding
var (
	moveDown    = fmt.Sprintf(RedColor, "j", "	move down")
	moveUp      = fmt.Sprintf(RedColor, "k", "	move up")
	moveLeft    = fmt.Sprintf(RedColor, "h", "	move left")
	moveRight   = fmt.Sprintf(RedColor, "l", "	move right")
	moveTop     = fmt.Sprintf(RedColor, "g", "	move top")
	moveBottom  = fmt.Sprintf(RedColor, "G", "	move bottom")
	pageDown    = fmt.Sprintf(RedColor, "ctrl-b", "page down")
	pageUp      = fmt.Sprintf(RedColor, "ctrl-f", "page up")
	stopApp     = fmt.Sprintf(RedColor, "ctrl-c", "stop tson")
	defaultNavi = strings.Join([]string{moveDown, moveUp, moveLeft,
		moveRight, moveTop, moveBottom, pageDown, pageUp, stopApp}, "\n")
)

// tree keybinding
var (
	hideNode           = fmt.Sprintf(RedColor, "h", "	hide children nodes")
	collaspeAllNode    = fmt.Sprintf(RedColor, "H", "	collaspe all nodes")
	expandNode         = fmt.Sprintf(RedColor, "l", "	expand children nodes")
	expandAllNode      = fmt.Sprintf(RedColor, "L", "	expand all children nodes")
	readFile           = fmt.Sprintf(RedColor, "r", "	read from file")
	saveFile           = fmt.Sprintf(RedColor, "s", "	save to file")
	addNewNode         = fmt.Sprintf(RedColor, "a", "	add new node")
	addNewValue        = fmt.Sprintf(RedColor, "A", "	add new value")
	clearChildrenNodes = fmt.Sprintf(RedColor, "d", "	clear children nodes")
	editNodes          = fmt.Sprintf(RedColor, "e", "	edit json with $EDITOR(only linux)")
	quitTson           = fmt.Sprintf(RedColor, "q", "	quit tson")
	editNodeValue      = fmt.Sprintf(RedColor, "Enter", "edit current node")
	searchNodes        = fmt.Sprintf(RedColor, "/", "	search nodes")
	toggleExpandNodes  = fmt.Sprintf(RedColor, "space", "	expand/collaspe nodes")
	moveNextParentNode = fmt.Sprintf(RedColor, "ctrl-j", "move to next parent node")
	movePreParentNode  = fmt.Sprintf(RedColor, "ctrl-k", "move to previous parent node")
	treeNavi           = strings.Join([]string{hideNode, collaspeAllNode, expandNode, expandAllNode,
		readFile, saveFile, addNewNode, addNewValue, clearChildrenNodes, editNodeValue, searchNodes,
		moveNextParentNode, movePreParentNode, editNodes, quitTson}, "\n")
)

type Navi struct {
	*tview.TextView
}

func NewNavi() *Navi {
	view := tview.NewTextView().SetDynamicColors(true)
	view.SetBorder(true).SetTitle("help").SetTitleAlign(tview.AlignLeft)
	navi := &Navi{TextView: view}
	return navi
}

func (n *Navi) UpdateView() {
	navi := strings.Join([]string{defaultNavi, "", treeNavi}, "\n")
	n.SetText(navi)
}

func (n *Navi) SetKeybindings(g *Gui) {
	n.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			g.Pages.HidePage(NaviPageName)
		}

		return event
	})
}
