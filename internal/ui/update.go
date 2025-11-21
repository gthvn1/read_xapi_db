package ui

import (
	"fmt"
	"sort"

	"github.com/rivo/tview"

	"example.com/readxapidb/internal/xapidb"
)

func UpdateStatus(tv *tview.TextView, n *xapidb.Node) {
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
