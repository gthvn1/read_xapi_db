package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"example.com/readxapidb/internal/args"
	"example.com/readxapidb/internal/fetch"
	"example.com/readxapidb/internal/theme"
	"example.com/readxapidb/internal/ui"
	"example.com/readxapidb/internal/xapidb"
)

func main() {
	args := args.GetArgs()

	var (
		data []byte
		err  error
	)

	data, err = fetch.DB(args)
	if err != nil {
		if args.Hostname == "" {
			fmt.Printf("failed to read %s: %s\n", args.FileName, err)
		} else {
			fmt.Printf("failed to fetch %s from %s: %s\n", args.FileName, args.Hostname, err)
		}
		os.Exit(1)
	}

	fmt.Printf("Read %d bytes from %s\n", len(data), args.FileName)

	DB, parse_err := xapidb.ParseXapiDB(data)
	if parse_err != nil && parse_err != io.EOF {
		fmt.Printf("failed to parse %s: %s\n", args.FileName, parse_err)
		os.Exit(1)
	}

	rootNode := DB.Root

	// Instead of printing the tree we will try to use the demo of navigable
	// tree view of current dir: https://github.com/rivo/tview/wiki/TreeView
	rootTree := ui.MakeTreeNode(rootNode)
	ui.LoadChildren(rootTree, rootNode)
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

		ui.UpdateStatus(status, node)

		// Load children if not already loaded
		if len(tn.GetChildren()) == 0 && len(node.Children) > 0 {
			ui.LoadChildren(tn, node)
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
