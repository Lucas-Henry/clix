package tui

import tea "github.com/charmbracelet/bubbletea"

type keyAction int

const (
	keyUp keyAction = iota
	keyDown
	keyLeft
	keyRight
	keyEnter
	keyTop
	keyBottom
	keyOpenWith
	keySearch
	keyContentSearch
	keyEsc
	keyQuit
	keyNone
)

func mapKey(msg tea.KeyMsg) keyAction {
	switch msg.String() {
	case "k", "up":
		return keyUp
	case "j", "down":
		return keyDown
	case "h", "left", "backspace":
		return keyLeft
	case "l", "right":
		return keyRight
	case "enter":
		return keyRight
	case "g":
		return keyTop
	case "G":
		return keyBottom
	case "o":
		return keyOpenWith
	case "/":
		return keySearch
	case "ctrl+f":
		return keyContentSearch
	case "esc":
		return keyEsc
	case "q", "ctrl+c":
		return keyQuit
	}
	return keyNone
}
