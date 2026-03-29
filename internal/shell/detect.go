package shell

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	ps "github.com/mitchellh/go-ps"
)

type Shell int

const (
	Unknown Shell = iota
	Bash
	Zsh
	Fish
	Sh
	Dash
	Ksh
)

func (s Shell) String() string {
	switch s {
	case Bash:
		return "bash"
	case Zsh:
		return "zsh"
	case Fish:
		return "fish"
	case Sh:
		return "sh"
	case Dash:
		return "dash"
	case Ksh:
		return "ksh"
	default:
		return "unknown"
	}
}

func Detect() Shell {
	if runtime.GOOS == "linux" {
		if s := fromProcFS(); s != Unknown {
			return s
		}
	}

	if s := fromParentProcess(); s != Unknown {
		return s
	}

	return fromEnv()
}

func fromProcFS() Shell {
	ppid := os.Getppid()
	commPath := filepath.Join("/proc", strconv.Itoa(ppid), "comm")
	data, err := os.ReadFile(commPath)
	if err != nil {
		return Unknown
	}
	return parse(strings.TrimSpace(string(data)))
}

func fromParentProcess() Shell {
	ppid := os.Getppid()
	proc, err := ps.FindProcess(ppid)
	if err != nil || proc == nil {
		return Unknown
	}
	return parse(filepath.Base(proc.Executable()))
}

func fromEnv() Shell {
	shellEnv := os.Getenv("SHELL")
	if shellEnv == "" {
		return Sh
	}
	return parse(filepath.Base(shellEnv))
}

func parse(name string) Shell {
	name = strings.ToLower(strings.TrimSpace(name))
	switch {
	case name == "bash" || strings.HasPrefix(name, "bash"):
		return Bash
	case name == "zsh" || strings.HasPrefix(name, "zsh"):
		return Zsh
	case name == "fish":
		return Fish
	case name == "dash":
		return Dash
	case name == "ksh" || name == "mksh" || name == "pdksh":
		return Ksh
	case name == "sh":
		return Sh
	default:
		return Unknown
	}
}
