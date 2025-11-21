package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"example.com/readxapidb/internal/xapidb"
)

func LoadChildren(tn *tview.TreeNode, n *xapidb.Node) {
	for _, c := range n.Children {
		tn.AddChild(MakeTreeNode(c))
	}
}

func MakeTreeNode(n *xapidb.Node) *tview.TreeNode {
	var label string

	// If there is a name attribute use it
	if name, ok := n.Attr["name"]; ok {
		label = fmt.Sprintf(" %s", name)
	} else {
		label = fmt.Sprintf(" %s", n.Name)
	}

	// If there is children print the number so you will know which
	// node can be unfold
	if len(n.Children) > 0 {
		label += fmt.Sprintf(" (%d)", len(n.Children))
	}

	// If there is a name__label add it, it not check if there is a ref.
	if nameLabel, ok := n.Attr["name__label"]; ok && len(nameLabel) > 0 {
		label += fmt.Sprintf(" [%s]", nameLabel)
	} else if ref, ok := n.Attr["ref"]; ok {
		label += fmt.Sprintf(" [%s]", ref)
	}

	tn := tview.NewTreeNode(label)
	tn.SetReference(n) // This maps the tree view with our node
	tn.SetSelectable(true)

	switch n.Name {
	case "database":
		tn.SetColor(tcell.ColorRed)
	case "table":
		tn.SetColor(tcell.ColorGreen)
	case "row":
		tn.SetColor(tcell.ColorBlue)
	default:
		tn.SetColor(tcell.ColorWhite)
	}

	// Just create the node, we will add children later
	return tn
}
