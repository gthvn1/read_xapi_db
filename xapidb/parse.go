package xapidb

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// Example of DB:
//
// <database>
//  <table name="Bond"/>
//  <table name="PCI">
//   <row ref="..." __ctime="..." ... />
//   <row ref="..." __ctime="..." ... />
//  </table>
//  <table name="VDI">
//   <row ref="..." ... />
//  </table>
// </database>

// XMLName  -> the element name (database, table, row)
// Attr     -> All attributes of the element (for table it is name))
// Children -> all sub-element (for database it is table)
//
// So for the example we have:
//
//     Node (database)
//     ├── XMLName = "database"
//     ├── Attr = []
//     ├── Children = [
//     │
//     │   Node (table)
//     │   ├── XMLName = "table"
//     │   ├── Attr = [ name="Bond" ]
//     │   ├── Children = []
//     │   └── InnerXML = ""
//     │
//     │   Node (table)
//     │   ├── XMLName = "table"
//     │   ├── Attr = [ name="PCI" ]
//     │   ├── Children = [
//     │   │
//     │   │   Node (row)
//     │   │   ├── XMLName = "row"
//     │   │   ├── Attr = [ ref="...", __ctime="...", ... ]
//     │   │   ├── Children = []
//     │   │   └── InnerXML = ""
//     │   │
//     │   │   Node (row)
//     │   │   ├── XMLName = "row"
//     │   │   ├── Attr = [ ref="...", __ctime="...", ... ]
//     │   │   ├── Children = []
//     │   │   └── InnerXML = ""
//     │   ]
//     │   └── InnerXML = "<row .../> <row .../>"
//     │
//     │   Node (table)
//     │   ├── XMLName = "table"
//     │   ├── Attr = [ name="VDI" ]
//     │   ├── Children = [
//     │   │   Node (row)
//     │   │   ├── XMLName = "row"
//     │   │   ├── Attr = [ ref="..." ]
//     │   │   ├── Children = []
//     │   │   └── InnerXML = ""
//     │   ]
//     │   └── InnerXML = "<row .../>"
//     ]
//     └── InnerXML = "<table .../> <table ...>...</table> ..."

// Using the parser we want to produce this tree:
//
//  Node(database)                           <-- root
//  ├─ Node(table)                            <-- table name="Bond"
//  │   Attr: { "name": "Bond" }
//  │   Children: []
//  │
//  ├─ Node(table)                            <-- table name="PCI"
//  │   Attr: { "name": "PCI" }
//  │   Children:
//  │     ├─ Node(row)
//  │     │   Attr: { "ref": "...", "__ctime": "...", ... }
//  │     │   Children: []
//  │     │
//  │     └─ Node(row)
//  │         Attr: { "ref": "...", "__ctime": "...", ... }
//  │         Children: []
//  │
//  └─ Node(table)                            <-- table name="VDI"
//      Attr: { "name": "VDI" }
//      Children:
//        └─ Node(row)
//            Attr: { "ref": "...", ... }
//            Children: []

type Node struct {
	Name     string
	Attr     map[string]string
	Children []*Node
	Parent   *Node // Will be usefull to deal with "cd .."
}

func PrintTree(n *Node) {
	var print func(node *Node, prefix string)

	print = func(node *Node, prefix string) {
		// Print node name
		fmt.Print(prefix)
		fmt.Print(node.Name)

		// Print attributes if any
		if len(node.Attr) > 0 {
			attrs := []string{}
			for k, v := range node.Attr {
				attrs = append(attrs, fmt.Sprintf(`%s="%s"`, k, v))
			}
			fmt.Print(" [", strings.Join(attrs, " "), "]")
		}
		fmt.Println()

		// Recurse into children
		for _, child := range node.Children {
			print(child, prefix+"  ") // indent 2 spaces per level
		}
	}

	print(n, "")
}

func ParseXapiDB(data []byte) (*Node, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))

	var stack []*Node
	var root *Node

	for {
		// Get the next XML token in the input stream
		tok, err := decoder.Token()
		if err != nil {
			// We reach the end of bytes
			if err == io.EOF {
				break
			}
			// Otherwise it is an error
			return nil, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			n := &Node{
				Name:     t.Name.Local,
				Attr:     map[string]string{},
				Children: []*Node{},
			}

			// Set attributes
			for _, a := range t.Attr {
				n.Attr[a.Name.Local] = a.Value
			}

			// Attach to parent if not root
			if len(stack) > 0 {
				parent := stack[len(stack)-1]
				n.Parent = parent
				parent.Children = append(parent.Children, n)
			} else {
				root = n
			}

			stack = append(stack, n)

		case xml.EndElement:
			stack = stack[:len(stack)-1]
		}
	}

	return root, nil
}
