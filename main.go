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
		searchHeight = 3
		debugHeight  = 5
		helpHeight   = 1
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
		AddItem(status, 0, 1, false)

	// We create 2 pages so we will be able to switch between
	// normal view and search view. Switch view is just normal view with
	// search bar on its top.
	normalLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainLayout, 0, 1, true).
		AddItem(debugView, debugHeight, 0, false).
		AddItem(help, helpHeight, 0, false)

	searchLayout := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(searchInput, searchHeight, 0, false). // keep search at the top
		AddItem(mainLayout, 0, 1, true).
		AddItem(debugView, debugHeight, 0, false).
		AddItem(help, helpHeight, 0, false)

	pages := tview.NewPages().
		AddPage("normal", normalLayout, true, true).
		AddPage("search", searchLayout, true, false)

	tview.Styles = theme.GruvboxDark

	app := tview.NewApplication()

	// Track which pane has focus
	var currentFocus tview.Primitive = tree

	// Set initial focus
	tree.SetBorderColor(tcell.ColorGreen)
	status.SetBorderColor(tcell.ColorWhite)

	// Set callbacks
	tree.SetSelectedFunc(ui.SelectedTreeCallback(status))
	status.SetSelectedFunc(ui.SelectedStatusCallback(status, debugView, app, tree, db))
	searchInput.SetDoneFunc(ui.DoneSearchCallback(app, tree, status, searchInput, debugView, db, pages))
	app.SetInputCapture(ui.InputCaptureCallback(app, tree, status, searchInput, pages, &currentFocus))

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
