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

	// Load children immediately for root
	loadChildren(rootTree, rootNode)

	rootTree.SetExpanded(true)

	tree := tview.NewTreeView().SetRoot(rootTree)

	// Set current node to first child if it exists
	if len(rootTree.GetChildren()) > 0 {
		tree.SetCurrentNode(rootTree.GetChildren()[0])
	} else {
		tree.SetCurrentNode(rootTree)
	}

	// Set border and title are done separatly otherwise the type of tree is
	// modified to tview.Box instead of TreeView !!!
	tree.SetBorder(true).SetTitle("XAPI DB")

	// We add a status view to print all row attributes for example
	status := tview.NewTextView().
		SetDynamicColors(true)
	status.SetBorder(true).SetTitle("Status")

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(tn *tview.TreeNode) {
		ref := tn.GetReference()
		if ref == nil {
			return // Selecting the root node does nothing.
		}

		node := ref.(*xapidb.Node)

		// Update status
		updateStatus(status, node)

		// Load children if not already loaded
		if len(tn.GetChildren()) == 0 && len(node.Children) > 0 {
			loadChildren(tn, node)
		}
		// Collapse if visible, expand if collapsed.
		tn.SetExpanded(!tn.IsExpanded())
	})

	layout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tree, 0, 1, true).
		AddItem(status, 40, 0, false)

	app := tview.NewApplication()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Detect "q"
		if event.Rune() == 'q' {
			app.Stop()
			return nil
		}

		return event
	})

	if err := app.SetRoot(layout, true).Run(); err != nil {
		panic(err)
	}
}

func loadChildren(tn *tview.TreeNode, n *xapidb.Node) {
	for _, c := range n.Children {
		tn.AddChild(makeTreeNode(c))
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

	if len(n.Children) > 0 {
		t.SetExpanded(false) // This should show expandable sign
	}

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

func updateStatus(tv *tview.TextView, n *xapidb.Node) {
	tv.Clear()

	fmt.Fprintf(tv, "[yellow]Name:[white] %s\n", n.Name)

	if len(n.Attr) > 0 {
		fmt.Fprintf(tv, "[yellow]Attributes:[white]\n")
		for k, v := range n.Attr {
			fmt.Fprintf(tv, "  %s = %q\n", k, v)
		}
	} else {
		fmt.Fprintf(tv, "[yellow]Attributes:[white] (none)\n")
	}

	fmt.Fprintf(tv, "[yellow]Children:[white] %d\n", len(n.Children))

	// Compute path
	path := ""
	cur := n
	for cur != nil {
		if cur.Parent != nil {
			path = "/" + cur.Name + path
		} else {
			path = "/database" + path
		}
		cur = cur.Parent
	}
	fmt.Fprintf(tv, "[yellow]Path:[white] %s\n", path)
}
