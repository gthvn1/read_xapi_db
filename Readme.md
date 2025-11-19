# XAPI DB CLI Explorer

A simple Go project to parse **XenServer XAPI XML databases** and navigate
them using [tview](https://github.com/rivo/tview).

---

## Features

- Parse XAPI XML database into a friendly tree structure.
- Navigate the tree 

## Installation

1. Clone the repository:

```bash
git clone https://github.com/gthvn1/read_xapi_db
cd read_xapi_db
```

2. Build the project:

```bash
go build .
```

3. Usage

```bash
./readxapidb database.xml
```
- where `database.xml` is the path of the XAPI database (currently it is hard coded to "xapi-db.ml").
