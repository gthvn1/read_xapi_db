package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"example.com/readxapidb/xapidb"
)

var XAPIDB = "./xapi-db.xml"

func main() {
	data, err := os.ReadFile(XAPIDB)
	if err != nil {
		fmt.Printf("failed to read %s: %s\n", XAPIDB, err)
		os.Exit(1)
	}

	fmt.Printf("Read %d bytes from %s\n", len(data), XAPIDB)

	rootNode, parse_err := xapidb.ParseXapiDB(data)
	if parse_err != nil && parse_err != io.EOF {
		fmt.Printf("failed to parse %s: %s\n", XAPIDB, parse_err)
		os.Exit(2)
	}

	//xapidb.PrintTree(t)

	// Instead of printing the tree we will try to use the demo of navigable
	// tree view of current dir: https://github.com/rivo/tview/wiki/TreeView
	rootTree := makeTreeNode(rootNode)
	tree := tview.NewTreeView().SetRoot(rootTree).SetCurrentNode(rootTree)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(tn *tview.TreeNode) {
		ref := tn.GetReference()
		if ref == nil {
			return // Selecting the root node does nothing.
		}

		node := ref.(*xapidb.Node)

		// Load children only once
		children := tn.GetChildren()
		if len(children) == 0 {
			for _, c := range node.Children {
				tn.AddChild(makeTreeNode(c))
			}
		}
		// Collapse if visible, expand if collapsed.
		tn.SetExpanded(!tn.IsExpanded())
	})

	if err := tview.NewApplication().SetRoot(tree, true).Run(); err != nil {
		panic(err)
	}
}

func makeTreeNode(n *xapidb.Node) *tview.TreeNode {
	label := n.Name

	// If there is a name attribute add it (it is the case for table)
	if name, ok := n.Attr["name"]; ok {
		label += fmt.Sprintf(" (%s)", name)
	}

	// If there is a ref add it (it is the case for row)
	if ref, ok := n.Attr["ref"]; ok {
		label += fmt.Sprintf(" [ref=%s]", ref)
	}

	t := tview.NewTreeNode(label)
	t.SetReference(n) // This maps the tree view with our node
	t.SetSelectable(true)

	switch n.Name {
	case "database":
		t.SetColor(tcell.ColorRed)
	case "table":
		t.SetColor(tcell.ColorGreen)
	case "row":
		t.SetColor(tcell.ColorBlue)
	default:
		t.SetColor(tcell.ColorWhite)
	}

	// Just create the node, we will add children later
	return t
}
