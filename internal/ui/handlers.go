package ui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"example.com/readxapidb/internal/xapidb"
)

// This function is called when the user selects tree
// by hitting Enter when selected
func SelectedTreeCallback(status *tview.Table) func(tn *tview.TreeNode) {
	return func(tn *tview.TreeNode) {
		// We are always setting a reference so let panic
		// if it is not the case...
		node := tn.GetReference().(*xapidb.Node)

		UpdateStatus(status, node)

		// Load children if not already loaded
		if len(tn.GetChildren()) == 0 && len(node.Children) > 0 {
			LoadChildren(tn, node)
			tn.SetExpanded(false)
		}

		// Collapse if visible, expand if collapsed.
		tn.SetExpanded(!tn.IsExpanded())
	}
}

// This function is called when the user selects status
// by hitting Enter when selected
func SelectedStatusCallback(
	status *tview.Table,
	debugView *tview.TextView,
	app *tview.Application,
	tree *tview.TreeView,
	db *xapidb.DB,
) func(row, column int) {
	return func(row, column int) {
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
			if retString := FollowOpaqueRef(app, tree, db, text); retString == "done" {
				fmt.Fprintf(debugView, "\n[green]Found the opaque reference")
			} else {
				fmt.Fprintf(debugView, "\n[red]%s", retString)
			}
		} else {
			fmt.Fprintf(debugView, "\n[blue]No match")
		}
	}
}

// handler which is called when the user is done entering text.
// The callback function is provided with the key that was pressed.
func DoneSearchCallback(
	app *tview.Application,
	tree *tview.TreeView,
	status *tview.Table,
	searchInput *tview.InputField,
	debugView *tview.TextView,
	db *xapidb.DB,
	pages *tview.Pages,
) func(key tcell.Key) {
	return func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			query := searchInput.GetText()

			if strings.HasPrefix(query, "OpaqueRef") {
				// Follow the reference
				result := FollowOpaqueRef(app, tree, db, query)
				debugView.Clear()
				fmt.Fprintf(debugView, "[yellow]Search:[white] %s\n", query)
				if result == "done" {
					// Update status
					if currentNode := tree.GetCurrentNode(); currentNode != nil {
						if ref := currentNode.GetReference(); ref != nil {
							node := ref.(*xapidb.Node)
							UpdateStatus(status, node)
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
			pages.SwitchToPage("normal")
			searchInput.SetText("")
			app.SetFocus(tree)
		}
	}
}

// InputCaptureCallback handles global keyboard input
func InputCaptureCallback(
	app *tview.Application,
	tree *tview.TreeView,
	status *tview.Table,
	searchInput *tview.InputField,
	pages *tview.Pages,
	currentFocus *tview.Primitive,
) func(event *tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		currentPage, _ := pages.GetFrontPage()
		inSearchMode := currentPage == "search"

		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case 'q':
				if *currentFocus != searchInput {
					app.Stop()
					return nil
				}

			case '/':
				// switch to search mode if not already in
				if !inSearchMode {
					pages.SwitchToPage("search")
					*currentFocus = searchInput
					app.SetFocus(*currentFocus)
				} else {
					pages.SwitchToPage("normal")
					*currentFocus = tree
					app.SetFocus(*currentFocus)
				}
				return nil

			case 'h', 'l':
				*currentFocus = ToggleFocus(app, currentFocus, tree, status)
				return nil
			}

		case tcell.KeyTab:
			if inSearchMode {
				*currentFocus = ToggleFocus(app, currentFocus, tree, status, searchInput)
			} else {
				*currentFocus = ToggleFocus(app, currentFocus, tree, status)
			}
			return nil
		}

		return event
	}
}
