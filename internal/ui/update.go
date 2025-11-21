package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/rivo/tview"

	"example.com/readxapidb/internal/xapidb"
)

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
