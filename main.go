package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

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

	_, parse_err := xapidb.ParseXapiDB(data)
	if parse_err != nil && parse_err != io.EOF {
		fmt.Println("failed to parse XAPI DB:", parse_err)
		os.Exit(2)
	}

	//xapidb.PrintTree(t)

	// Instead of printing the tree we will try to use the demo of navigable
	// tree view of current dir: https://github.com/rivo/tview/wiki/TreeView
	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(file.IsDir())
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	add(root, rootDir)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)
			add(node, path)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	if err := tview.NewApplication().SetRoot(tree, true).Run(); err != nil {
		panic(err)
	}
}
