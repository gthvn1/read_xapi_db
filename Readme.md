# XAPI DB CLI Explorer

A Go-based interactive viewer for **XCP-ng / XenServer XAPI XML databases**, displayed
using [tview](https://github.com/rivo/tview).

This tool lets you **navigate, search, and inspect** the XAPI DB just like a filesystem
tree. It can load a database **locally** or **directly from a remote XCP-ng host over
SSH/SFTP**.

---

## Features

- Parse the full XAPI XML database into a navigable tree.
- Browse tables and rows interactively (expand/collapse).
- View attributes sorted alphabetically.
- **NEW:** Fetch the XAPI DB directly from a remote XCP-ng host via SSH/SFTP.
- **WIP:** Follow cross-references (`OpaqueRef:*`) between rows.

## Installation

### Clone the repository:

```bash
git clone https://github.com/gthvn1/read_xapi_db
cd read_xapi_db
```

### Build the project:

```bash
go build .
```

### Usage

#### Local file mode

```bash
./readxapidb --file ./xapi-db.xml
```

#### Remote mode (NEW)

You can fetch the database directly from any XenServer/XCP-ng host:
```bash
./readxapidb \
    --hostname xenhost \
    --username root \
    --password mypassword \
    --file /var/lib/xcp/state.db
```
- Arguments

| Flag         | Description                                           |
| ------------ | ----------------------------------------------------- |
| `--file`     | Path to the database (local OR remote). **Required.** |
| `--hostname` | Remote hostname or IP. Leave empty to use local mode. |
| `--username` | SSH username (remote mode only).                      |
| `--password` | SSH password (remote mode only).                      |

If `--hostname` is not provided, the tool loads the file locally.

---

<img src="https://github.com/gthvn1/read_xapi_db/blob/master/screenshot.png">
