# XAPI DB CLI Explorer

A simple Go project to parse **XenServer XAPI XML databases** and navigate
them using a command-line interface (CLI).

The project will support **listing tables and rows** (`ls`) and
**changing context** (`cd`) within the XML tree. Future improvements may
include `find` and more advanced queries.

---

## Features

- Parse XAPI XML database into a friendly tree structure.
- Navigate the tree with familiar CLI commands:
  - `ls` — list tables or rows in the current node.
  - `cd <table/row>` — move into a table or back with `cd ..`.
- Display node attributes conveniently.
- Lightweight and easy to extend.

---

## Installation

1. Clone the repository:

```bash
git clone https://github.com/gthvn1/read_xapi_db
cd read_xapi_db
```

2. Build the project:

```bash
go build -o xdb-cli .
```

3. Usage

```bash
./xdb-cli database.xml
```
- where `database.xml` is the path of the XAPI database.
