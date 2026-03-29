package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Lucas-Henry/clix/internal/ipc"
	"github.com/Lucas-Henry/clix/internal/tui"
)

var version = "0.1.0"

func main() {
	var (
		startDir   = flag.String("dir", "", "starting directory (default: cwd)")
		noIPC      = flag.Bool("no-ipc", false, "disable unix socket IPC server")
		showVer    = flag.Bool("version", false, "print version and exit")
		socketPath = flag.Bool("socket-path", false, "print socket path and exit")
	)
	flag.Parse()

	if *showVer {
		fmt.Println("clix", version)
		os.Exit(0)
	}

	dir := *startDir
	if dir == "" && flag.NArg() > 0 {
		dir = flag.Arg(0)
	}

	if *socketPath {
		fmt.Println(ipc.SocketPath())
		os.Exit(0)
	}

	model := tui.New(dir, !*noIPC)

	var srv *ipc.Server
	if !*noIPC {
		srv = ipc.NewServer(
			func() (string, string) {
				return model.CWD(), model.SelectedName()
			},
			func(path string) error {
				return model.Navigate(path)
			},
		)
		if err := srv.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "clix: ipc server failed: %v\n", err)
		} else {
			defer srv.Stop()
		}
	}

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "clix: %v\n", err)
		os.Exit(1)
	}
}
