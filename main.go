package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

var XAPIDB = "./xapi-db.xml"

func checkOrDie(e error) {
	if e != nil {
		panic(e)
	}
}

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

type Node struct {
	XMLName  xml.Name
	Attr     []xml.Attr `xml:",any,attr"`
	Children []Node     `xml:",any"`
	InnerXML string     `xml:",innerxml"`
}

func main() {
	data, err := os.ReadFile(XAPIDB)
	checkOrDie(err)
	fmt.Printf("Read %d bytes from %s\n", len(data), XAPIDB)

	var root Node
	checkOrDie(xml.Unmarshal(data, &root))

	fmt.Println("Root element: ", root.XMLName.Local)

	for _, table := range root.Children {
		if table.XMLName.Local != "table" {
			continue
		}

		var tableName string
		for _, a := range table.Attr {
			if a.Name.Local == "name" {
				tableName = a.Value
			}
		}

		for _, row := range table.Children {
			if row.XMLName.Local != "row" {
				continue
			}

			fmt.Printf("Row in table %s:\n", tableName)
			for _, a := range row.Attr {
				fmt.Printf("  %s = %s\n", a.Name.Local, a.Value)
			}
		}
	}

}
