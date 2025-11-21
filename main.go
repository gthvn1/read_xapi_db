package main

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"example.com/readxapidb/theme"
	"example.com/readxapidb/xapidb"
)

func main() {
	dbFile := "./xapi-db.xml"

	switch len(os.Args) {
	case 2:
		dbFile = os.Args[1]
	case 1:
	// Use default
	default:
		fmt.Printf("Usage: %s [dbfile.xml]\n", os.Args[0])
		fmt.Printf("If no file is provided, the default is: %s\n", dbFile)
		os.Exit(1)
	}

	data, err := os.ReadFile(dbFile)
	if err != nil {
		fmt.Printf("failed to read %s: %s\n", dbFile, err)
		os.Exit(1)
	}

	fmt.Printf("Read %d bytes from %s\n", len(data), dbFile)

	DB, parse_err := xapidb.ParseXapiDB(data)
	if parse_err != nil && parse_err != io.EOF {
		fmt.Printf("failed to parse %s: %s\n", dbFile, parse_err)
		os.Exit(2)
	}

	rootNode := DB.Root

	// Instead of printing the tree we will try to use the demo of navigable
	// tree view of current dir: https://github.com/rivo/tview/wiki/TreeView
	rootTree := makeTreeNode(rootNode)
	loadChildren(rootTree, rootNode)
	rootTree.SetExpanded(true)

	tree := tview.NewTreeView()

	// Set current node to first child if it exists
	if len(rootTree.GetChildren()) > 0 {
		tree.SetCurrentNode(rootTree.GetChildren()[0])
	} else {
		tree.SetCurrentNode(rootTree)
	}

	// Set border and title are done separatly otherwise the type of tree is
	// modified to tview.Box instead of TreeView !!!
	tree.SetRoot(rootTree).
		SetBorder(true).
		SetTitle("XAPI DB")

	// We add a status view to print all row attributes for example
	status := tview.NewTextView()
	status.SetDynamicColors(true).
		SetScrollable(true).
		SetBorder(true).
		SetTitle("Status")

	// Add search input (initially hidden)
	searchInput := tview.NewInputField()
	searchInput.SetLabel("Seach: ").
		SetFieldWidth(50).
		SetBorder(true).
		SetTitle("Search (ESC to cancel)")

	// Search results view
	searchResults := tview.NewTextView()
	searchResults.SetDynamicColors(true).
		SetScrollable(true).
		SetBorder(true).
		SetTitle("Search Results")

	// Track which pane has focus
	var currentFocus tview.Primitive = tree
	searchMode := false

	// This function is called when the user selects this
	// node by hitting Enter when selected
	tree.SetSelectedFunc(func(tn *tview.TreeNode) {
		// We are always setting a reference so let panic
		// if it is not the case...
		node := tn.GetReference().(*xapidb.Node)

		updateStatus(status, node)

		// Load children if not already loaded
		if len(tn.GetChildren()) == 0 && len(node.Children) > 0 {
			loadChildren(tn, node)
			tn.SetExpanded(false)
		}

		// Collapse if visible, expand if collapsed.
		tn.SetExpanded(!tn.IsExpanded())
	})

	// Add help footer
	help := tview.NewTextView()
	help.SetTextAlign(tview.AlignCenter).SetDynamicColors(true)
	help.SetText("[yellow]'q'[white]=quit | [yellow]'/'[white]=search | [yellow]'Space/Enter'[white]=expand/collapse")
	help.SetBackgroundColor(tcell.ColorDefault)

	// Create main Layout with tree and status
	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tree, 0, 1, true).
		AddItem(status, 0, 1, false)

	// We create 2 pages so we will be able to switch between
	// normal view and search view
	normalLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainLayout, 0, 1, true).
		AddItem(help, 1, 0, false)

	searchLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(searchInput, 3, 0, true).
		AddItem(searchResults, 0, 1, true)

	pages := tview.NewPages().
		AddPage("normal", normalLayout, true, true).
		AddPage("search", searchLayout, true, false)

	tview.Styles = theme.GruvboxDark

	app := tview.NewApplication()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				if !searchMode {
					app.Stop()
					return nil
				}

			case '/':
				// switch to search mode if not already in
				if !searchMode {
					searchMode = true
					pages.SwitchToPage("search")
					return nil
				}

			case 'h', 'l':
				// switch between tree and status
				if currentFocus == tree {
					currentFocus = status
					tree.SetBorderColor(tcell.ColorWhite)
					status.SetBorderColor(tcell.ColorGreen)
				} else {
					currentFocus = tree
					tree.SetBorderColor(tcell.ColorGreen)
					status.SetBorderColor(tcell.ColorWhite)
				}
				app.SetFocus(currentFocus)
				return nil
			}

		case tcell.KeyTab:
			// switch between tree and status
			if currentFocus == tree {
				currentFocus = status
				tree.SetBorderColor(tcell.ColorWhite)
				status.SetBorderColor(tcell.ColorGreen)
			} else {
				currentFocus = tree
				tree.SetBorderColor(tcell.ColorGreen)
				status.SetBorderColor(tcell.ColorWhite)
			}
			app.SetFocus(currentFocus)
			return nil

		case tcell.KeyEscape:
			searchMode = false
			pages.SwitchToPage("normal")
			currentFocus = tree
			tree.SetBorderColor(tcell.ColorGreen)
			app.SetFocus(tree)
			return nil
		}

		return event
	})

	// Set initial focus
	tree.SetBorderColor(tcell.ColorGreen)
	status.SetBorderColor(tcell.ColorWhite)

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}

func loadChildren(tn *tview.TreeNode, n *xapidb.Node) {
	for _, c := range n.Children {
		tn.AddChild(makeTreeNode(c))
	}
}

func makeTreeNode(n *xapidb.Node) *tview.TreeNode {
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

func updateStatus(tv *tview.TextView, n *xapidb.Node) {
	tv.Clear()

	fmt.Fprintf(tv, "[yellow]Name:[white] %s\n", n.Name)

	if len(n.Attr) > 0 {
		fmt.Fprintf(tv, "[yellow]Attributes:[white]\n")
		// We first sort keys
		keys := make([]string, 0, len(n.Attr))
		for k := range n.Attr {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			fmt.Fprintf(tv, "  [orange]%s[white] = %q\n", k, n.Attr[k])
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
