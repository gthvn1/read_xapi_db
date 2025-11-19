package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// DataNode represents your XAPI data structure
type DataNode struct {
	Name     string
	Children []*DataNode
	Content  string
	Data     map[string]string
}

func main() {
	// Create sample tree (replace this with your XML parser)
	root := &DataNode{
		Name:    "XAPI Database",
		Content: "Root of XAPI database",
		Data: map[string]string{
			"Version": "1.0",
			"Host":    "xenserver-01",
		},
		Children: []*DataNode{
			{
				Name:    "VMs",
				Content: "Virtual Machines Collection",
				Data: map[string]string{
					"Total": "2",
				},
				Children: []*DataNode{
					{
						Name:    "vm-001",
						Content: "Ubuntu Server VM",
						Data: map[string]string{
							"OS":     "Ubuntu 22.04",
							"CPUs":   "4",
							"RAM":    "8GB",
							"Status": "Running",
						},
					},
					{
						Name:    "vm-002",
						Content: "Windows 10 VM",
						Data: map[string]string{
							"OS":     "Windows 10",
							"CPUs":   "2",
							"RAM":    "4GB",
							"Status": "Stopped",
						},
					},
				},
			},
			{
				Name:    "Networks",
				Content: "Network Configurations",
				Data: map[string]string{
					"Total": "2",
				},
				Children: []*DataNode{
					{
						Name:    "net-001",
						Content: "Bridge Network",
						Data: map[string]string{
							"Type":   "Bridge",
							"Bridge": "xenbr0",
							"MTU":    "1500",
						},
					},
					{
						Name:    "net-002",
						Content: "NAT Network",
						Data: map[string]string{
							"Type":   "NAT",
							"Subnet": "192.168.100.0/24",
						},
					},
				},
			},
			{
				Name:    "Storage",
				Content: "Storage Repositories",
				Data: map[string]string{
					"Total": "1",
				},
				Children: []*DataNode{
					{
						Name:    "sr-001",
						Content: "Local Storage",
						Data: map[string]string{
							"Type": "LVM",
							"Size": "500GB",
							"Used": "120GB",
						},
					},
				},
			},
		},
	}

	app := tview.NewApplication()

	// Detail view for selected node
	detailView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)

	// Helper function to format node details
	formatDetails := func(node *DataNode) string {
		text := fmt.Sprintf("[yellow]%s[white]\n\n", node.Name)
		text += fmt.Sprintf("%s\n\n", node.Content)

		if len(node.Data) > 0 {
			text += "[cyan]Details:[white]\n"
			for key, value := range node.Data {
				text += fmt.Sprintf("  [green]%s:[white] %s\n", key, value)
			}
		}

		if len(node.Children) > 0 {
			text += fmt.Sprintf("\n[cyan]Children:[white] %d\n", len(node.Children))
		}

		return text
	}

	// Create tree view
	treeView := tview.NewTreeView()

	// Build tree recursively
	var addNode func(*tview.TreeNode, *DataNode)
	addNode = func(target *tview.TreeNode, data *DataNode) {
		for _, child := range data.Children {
			// Add expand/collapse indicator
			prefix := ""
			if len(child.Children) > 0 {
				prefix = "[+] "
			}

			node := tview.NewTreeNode(prefix + child.Name).
				SetReference(child).
				SetSelectable(true).
				SetExpanded(false) // Start collapsed

			if len(child.Children) > 0 {
				node.SetColor(tview.Styles.SecondaryTextColor)
			}

			target.AddChild(node)
			addNode(node, child)
		}
	}

	// Create root node
	rootNode := tview.NewTreeNode(root.Name).
		SetReference(root).
		SetSelectable(true).
		SetExpanded(false). // Start collapsed
		SetColor(tview.Styles.PrimaryTextColor)

	treeView.SetRoot(rootNode).
		SetCurrentNode(rootNode)

	// Add children to root
	addNode(rootNode, root)

	// Update detail view when selection changes
	treeView.SetChangedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref != nil {
			dataNode := ref.(*DataNode)
			detailView.SetText(formatDetails(dataNode))
		}
	})

	// Handle expand/collapse to update indicators
	treeView.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			return
		}

		dataNode := ref.(*DataNode)
		if len(dataNode.Children) == 0 {
			return // Leaf node, nothing to toggle
		}

		// Toggle expansion
		node.SetExpanded(!node.IsExpanded())

		// Update indicator
		if node.IsExpanded() {
			node.SetText("[-] " + dataNode.Name)
		} else {
			node.SetText("[+] " + dataNode.Name)
		}
	})

	// Initial detail display
	detailView.SetText(formatDetails(root))

	// Help text
	helpView := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[yellow]Navigation:[white] ↑↓ Move | Enter Expand/Collapse | [yellow]Quit:[white] Ctrl+C or q")

	// Create layout
	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(treeView, 0, 1, true).
			AddItem(helpView, 1, 0, false), 0, 1, true).
		AddItem(detailView, 0, 1, false)

	// Add keyboard shortcuts
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'q' {
			app.Stop()
			return nil
		}
		return event
	})

	if err := app.SetRoot(flex, true).SetFocus(treeView).Run(); err != nil {
		panic(err)
	}
}
