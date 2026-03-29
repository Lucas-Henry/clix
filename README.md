# clix

A keyboard-driven terminal file explorer with vim keybindings, file preview, and a Unix socket IPC interface for AI agent integration.

## Installation

### From source

```sh
git clone https://github.com/Lucas-Henry/clix.git
cd clix
go install ./cmd/...
```

The `clix` binary will be placed in `$GOPATH/bin` (or `$HOME/go/bin`). Make sure that directory is in your `$PATH`.

### Quick install

```sh
go install github.com/Lucas-Henry/clix/cmd@latest
```

## Usage

```sh
clix [flags] [directory]
```

| Flag | Description |
|------|-------------|
| `--dir <path>` | Start in the given directory (default: cwd) |
| `--no-ipc` | Disable the Unix socket IPC server |
| `--socket-path` | Print the socket path and exit |
| `--version` | Print version and exit |

## Keybindings

| Key | Action |
|-----|--------|
| `j` / `down` | Move cursor down |
| `k` / `up` | Move cursor up |
| `h` / `backspace` | Go to parent directory |
| `l` / `enter` | Enter directory or focus preview |
| `gg` | Jump to first entry |
| `G` | Jump to last entry |
| `o` | Open with (nano / vim / nvim / cat / custom) |
| `q` / `ctrl+c` | Quit |

When the preview pane is focused, `j`/`k` scroll the file content.

## IPC (Unix Socket)

When clix starts, it creates a Unix socket at `/tmp/clix-<pid>.sock`. External processes (AI agents, MCP servers) can connect and send newline-delimited JSON commands.

### Get socket path

```sh
clix --socket-path
# /tmp/clix-12345.sock
```

Or read it from the running instance:

```sh
SOCK=$(clix --socket-path)
```

### Protocol

**Request**

```json
{ "action": "get_state" }
{ "action": "navigate", "path": "/home/user/projects" }
```

**Response**

```json
{ "ok": true, "cwd": "/home/user/projects", "selected": "clix" }
{ "ok": false, "error": "path required" }
```

### Actions

| Action | Fields | Description |
|--------|--------|-------------|
| `get_state` | — | Returns current working directory and selected entry name |
| `navigate` | `path` (string) | Navigates the UI to the given absolute path |

### Example: connect from shell

```sh
SOCK=$(clix --socket-path)
echo '{"action":"navigate","path":"/tmp"}' | nc -U "$SOCK"
```

### Example: connect from Python agent

```python
import socket, json

sock_path = "/tmp/clix-12345.sock"

with socket.socket(socket.AF_UNIX, socket.SOCK_STREAM) as s:
    s.connect(sock_path)
    s.sendall(json.dumps({"action": "get_state"}).encode() + b"\n")
    response = json.loads(s.recv(4096))
    print(response)
```

## Shell detection

clix automatically detects the active shell using the following strategy:

1. Reads `/proc/<ppid>/comm` (Linux only)
2. Inspects the parent process name via the OS process list (Linux + macOS)
3. Falls back to `$SHELL` environment variable

Supported shells: `bash`, `zsh`, `fish`, `dash`, `ksh`, `sh`.

## Project structure

```
clix/
├── cmd/
│   └── main.go              # entrypoint, flags, program boot
├── internal/
│   ├── tui/
│   │   ├── app.go           # bubbletea root model
│   │   ├── pane_tree.go     # left pane: file tree navigation
│   │   ├── pane_preview.go  # right pane: file content preview
│   │   ├── keymap.go        # vim keybinding mapping
│   │   └── styles.go        # lipgloss style definitions
│   ├── fs/
│   │   ├── walker.go        # directory reading, entry sorting
│   │   └── preview.go       # file preview, binary detection
│   ├── shell/
│   │   └── detect.go        # shell detection cascade
│   ├── opener/
│   │   └── opener.go        # open-with command builder
│   └── ipc/
│       └── server.go        # unix socket IPC server
└── go.mod
```

## License

MIT
