package main

import (
	"fmt"
	"io"
	"os"
	"strings"

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

	data, err := fetch.DB(args)
	if err != nil {
		if args.Hostname == "" {
			fmt.Printf("failed to read %s: %s\n", args.FileName, err)
		} else {
			fmt.Printf("failed to fetch %s from %s: %s\n", args.FileName, args.Hostname, err)
		}
		os.Exit(1)
	}

	fmt.Printf("Read %d bytes from %s\n", len(data), args.FileName)

	db, err := xapidb.ParseXapiDB(data)
	if err != nil && err != io.EOF {
		fmt.Printf("failed to parse %s: %s\n", args.FileName, err)
		os.Exit(1)
	}

	rootNode := db.Root

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

	const (
		statusPaneWidth = 75
		searchHeight    = 3
		debugHeight     = 5
		helpHeight      = 1
	)

	// Set border and title are done separatly otherwise the type of tree is
	// modified to tview.Box instead of TreeView !!!
	tree.SetRoot(rootTree).
		SetBorder(true).
		SetTitle("XAPI DB")

	// We add a status view to print all row attributes for example
	status := tview.NewTable()
	status.SetBorders(false).
		SetSelectable(true, false).
		SetSelectedStyle(tcell.Style{}.
			Background(tcell.NewHexColor(0x504945)).
			Foreground(tcell.NewHexColor(0xfabd2f))).
		SetBorder(true).
		SetTitle("Attributes")

	// Add search input (initially hidden)
	searchInput := tview.NewInputField()
	searchInput.SetLabel("Seach: ").
		SetBorder(true).
		SetTitle("Search")

	// Create a debug/info view
	debugView := tview.NewTextView()
	debugView.SetDynamicColors(true).
		SetBorder(true).
		SetTitle("Debug")

	// Add help footer
	help := tview.NewTextView()
	help.SetTextAlign(tview.AlignCenter).SetDynamicColors(true)
	help.SetText("[yellow]'q'[white]=quit | [yellow]'/'[white]=search | [yellow]'Space/Enter'[white]=expand/collapse")
	help.SetBackgroundColor(tcell.ColorDefault)

	// Create main Layout with tree and status
	mainLayout := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(tree, 0, 1, true).
		AddItem(status, statusPaneWidth, 0, false)

	// We create 2 pages so we will be able to switch between
	// normal view and search view (TODO: search view has been removed)
	normalLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainLayout, 0, 1, true).
		AddItem(debugView, debugHeight, 0, false).
		AddItem(help, helpHeight, 0, false)

	pages := tview.NewPages().
		AddPage("normal", normalLayout, true, true)

	tview.Styles = theme.GruvboxDark

	app := tview.NewApplication()

	// Track which pane has focus
	var currentFocus tview.Primitive = tree
	searchMode := false

	// Set initial focus
	tree.SetBorderColor(tcell.ColorGreen)
	status.SetBorderColor(tcell.ColorWhite)

	// -------------------------------------------------------------------
	// Callbacks
	// -------------------------------------------------------------------

	// This function is called when the user selects status
	// by hitting Enter when selected
	selectedStatusCallback := func(row, column int) {
		valueCell := status.GetCell(row, 1)
		if valueCell == nil {
			return
		}

		text := valueCell.Text

		// Has we have color on OpaqueRef we use the reference to get the
		// raw string (it has been set during update). If there is no ref
		// keep using the text.
		if ref := valueCell.GetReference(); ref != nil {
			text = ref.(string)
		}

		debugView.Clear()
		fmt.Fprintf(debugView, "Text: %s (len=%d)", text, len(text))
		if len(text) > 0 {
			preview := text[:min(3, len(text))]
			fmt.Fprintf(debugView, "\nFirst 3 chars: %q", preview)
		}

		if strings.HasPrefix(text, "OpaqueRef") {
			if retString := ui.FollowOpaqueRef(app, tree, db, text); retString == "done" {
				fmt.Fprintf(debugView, "\n[green]Found the opaque reference")
			} else {
				fmt.Fprintf(debugView, "\n[red]%s", retString)
			}
		} else {
			fmt.Fprintf(debugView, "\n[blue]No match")
		}
	}

	// This function is called when the user selects tree
	// by hitting Enter when selected
	selectedTreeCallback := func(tn *tview.TreeNode) {
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
	}

	// handler which is called when the user is done entering text.
	// The callback function is provided with the key that was pressed,
	doneSearchCallback := func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			query := searchInput.GetText()

			if strings.HasPrefix(query, "OpaqueRef") {
				// Follow the reference
				result := ui.FollowOpaqueRef(app, tree, db, query)
				debugView.Clear()
				fmt.Fprintf(debugView, "[yellow]Search:[white] %s\n", query)
				if result == "done" {
					// Update status
					if currentNode := tree.GetCurrentNode(); currentNode != nil {
						if ref := currentNode.GetReference(); ref != nil {
							node := ref.(*xapidb.Node)
							ui.UpdateStatus(status, node)
						}
					}
					fmt.Fprintf(debugView, "[green]Found reference!")
				} else {
					fmt.Fprintf(debugView, "[red]%s", result)
				}
			} else {
				// TODO: General text search in nodes
				debugView.Clear()
				fmt.Fprintf(debugView, "[yellow]Searching for:[white] %s\n", query)
				fmt.Fprintf(debugView, "[blue]Text search not implemented yet")
			}

		case tcell.KeyEscape:
			// Cancel search
			searchMode = false
			normalLayout.RemoveItem(searchInput)
			searchInput.SetText("")
			app.SetFocus(tree)
		}
	}

	tree.SetSelectedFunc(selectedTreeCallback)
	status.SetSelectedFunc(selectedStatusCallback)
	searchInput.SetDoneFunc(doneSearchCallback)

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
					// Keep the help at the end
					normalLayout.RemoveItem(help)
					normalLayout.
						AddItem(searchInput, searchHeight, 0, false).
						AddItem(help, helpHeight, 0, false)
					currentFocus = searchInput
					app.SetFocus(currentFocus)
				} else {
					searchMode = false
					normalLayout.RemoveItem(searchInput)
					currentFocus = tree
					app.SetFocus(currentFocus)
				}
				return nil

			case 'h', 'l':
				currentFocus = ui.ToggleFocus(app, &currentFocus, tree, status)
				return nil
			}

		case tcell.KeyTab:
			if searchMode {
				currentFocus = ui.ToggleFocus(app, &currentFocus, tree, status, searchInput)
			} else {
				currentFocus = ui.ToggleFocus(app, &currentFocus, tree, status)
			}
			return nil
		}

		return event
	})

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
