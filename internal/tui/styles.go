package tui

import "github.com/charmbracelet/lipgloss"

var (
	colorBorder    = lipgloss.Color("240")
	colorBorderFoc = lipgloss.Color("39")
	colorSelected  = lipgloss.Color("39")
	colorDim       = lipgloss.Color("243")
	colorDir       = lipgloss.Color("33")
	colorFile      = lipgloss.Color("252")
	colorSymlink   = lipgloss.Color("43")
	colorHeader    = lipgloss.Color("252")
	colorStatus    = lipgloss.Color("243")
	colorIPCActive = lipgloss.Color("40")
	colorIPCOff    = lipgloss.Color("240")
	colorBinary    = lipgloss.Color("196")
	colorLineNum   = lipgloss.Color("240")

	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorHeader)

	styleDim = lipgloss.NewStyle().
			Foreground(colorDim)

	styleDir = lipgloss.NewStyle().
			Foreground(colorDir).
			Bold(true)

	styleFile = lipgloss.NewStyle().
			Foreground(colorFile)

	styleSymlink = lipgloss.NewStyle().
			Foreground(colorSymlink)

	styleSelectedItem = lipgloss.NewStyle().
				Background(colorSelected).
				Foreground(lipgloss.Color("0")).
				Bold(true)

	styleLineNum = lipgloss.NewStyle().
			Foreground(colorLineNum)

	styleBinaryWarn = lipgloss.NewStyle().
			Foreground(colorBinary).
			Bold(true)

	panelBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder)

	panelBorderFocused = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorBorderFoc)
)
