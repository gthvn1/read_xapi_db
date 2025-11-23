package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"example.com/readxapidb/internal/xapidb"
)

func ToggleFocus(app *tview.Application, current *tview.Primitive, items ...tview.Primitive) tview.Primitive {
	// Find current
	currentIdx := -1
	for i, item := range items {
		if *current == item {
			currentIdx = i
			break
		}
	}

	nextIdx := (currentIdx + 1) % len(items)

	// Update all borders
	for i, item := range items {
		if boxItem, ok := item.(interface{ SetBorderColor(tcell.Color) *tview.Box }); ok {
			if i == nextIdx {
				boxItem.SetBorderColor(tcell.ColorGreen)
			} else {
				boxItem.SetBorderColor(tcell.ColorWhite)
			}
		}
	}

	*current = items[nextIdx]
	app.SetFocus(*current)

	return *current
}

func UpdateStatus(tv *tview.Table, n *xapidb.Node) {
	tv.Clear()

	row := 0

	// Name
	tv.SetCell(row, 0, tview.NewTableCell("[yellow]Name[white]"))
	tv.SetCell(row, 1, tview.NewTableCell(n.Name))
	row++

	// Attributes
	tv.SetCell(row, 0, tview.NewTableCell("[yellow]Attributes[white]"))
	row++

	if len(n.Attr) > 0 {
		// We first sort keys
		keys := make([]string, 0, len(n.Attr))
		for k := range n.Attr {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := n.Attr[k]
			keyCell := tview.NewTableCell("  [orange]" + k + "[white]")
			valCell := tview.NewTableCell(v)

			// Highlight OpaqueRefs that we will able to follow (WIP)
			if strings.HasPrefix(v, "OpaqueRef:") {
				valCell = tview.NewTableCell("[blue]" + v + "[white]")
				valCell.SetReference(v) // Store the raw string to be able to follow the OpaqueRef
				valCell.SetSelectable(true)
			}

			tv.SetCell(row, 0, keyCell)
			tv.SetCell(row, 1, valCell)

			row++
		}
	} else {
		tv.SetCell(row, 0, tview.NewTableCell("[yellow]Attributes[white]"))
		tv.SetCell(row, 1, tview.NewTableCell("(none)"))

		row++
	}

	// Children count
	tv.SetCell(row, 0, tview.NewTableCell("[yellow]Children[white]"))
	tv.SetCell(row, 1, tview.NewTableCell(fmt.Sprintf("%d", len(n.Children))))
	row++

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

	tv.SetCell(row, 0, tview.NewTableCell("[yellow]Path[white]"))
	tv.SetCell(row, 1, tview.NewTableCell(path))
}

func FollowOpaqueRef(app *tview.Application, tree *tview.TreeView, DB *xapidb.DB, ref string) string {
	// Find node using the DB ref index
	target, ok := DB.RefIndex[ref]
	if !ok {
		return fmt.Sprintf("Failed to find %s in RefIndex", ref)
	}

	// Target parent is always the table
	table := target.Parent
	if table == nil {
		return "failed to find the table parent"
	}

	root := tree.GetRoot()
	root.SetExpanded(true)

	// Find table node inside the tree
	var tableTreeNode *tview.TreeNode
	for _, tn := range root.GetChildren() {
		if tn.GetReference() == table {
			tableTreeNode = tn
			break
		}
	}

	if tableTreeNode == nil {
		return "failed to find the corresponding table in TreeView"
	}

	// Load rows of the table if not loaded yet
	LoadChildren(tableTreeNode, table)
	tableTreeNode.SetExpanded(true)

	// Now find the row node inside the table
	for _, rowNode := range tableTreeNode.GetChildren() {
		if rowNode.GetReference() == target {
			// scroll to it and select it
			rowNode.SetExpanded(true)
			tree.SetCurrentNode(rowNode)
			app.SetFocus(tree)
			return "done"
		}
	}

	return "failed to find the row node inside the table"
}
