package opener

import (
	"os/exec"
	"path/filepath"
)

type Editor int

const (
	Nano Editor = iota
	Vim
	Nvim
	Cat
	Custom
)

var EditorNames = []string{"nano", "vim", "nvim", "cat"}

func (e Editor) String() string {
	switch e {
	case Nano:
		return "nano"
	case Vim:
		return "vim"
	case Nvim:
		return "nvim"
	case Cat:
		return "cat"
	default:
		return "custom"
	}
}

func Build(editor Editor, customCmd string, filePath string) *exec.Cmd {
	dir := filepath.Dir(filePath)

	var cmd *exec.Cmd
	if editor == Custom {
		cmd = exec.Command(customCmd, filePath)
	} else {
		cmd = exec.Command(editor.String(), filePath)
	}

	cmd.Dir = dir
	return cmd
}

func Available(editor Editor) bool {
	_, err := exec.LookPath(editor.String())
	return err == nil
}
