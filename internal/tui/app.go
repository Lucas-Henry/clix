package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	clixfs "github.com/Lucas-Henry/clix/internal/fs"
	"github.com/Lucas-Henry/clix/internal/opener"
	"github.com/Lucas-Henry/clix/internal/shell"
)

type focusPane int

const (
	focusTree focusPane = iota
	focusPreview
)

type openWithState int

const (
	openWithOff openWithState = iota
	openWithActive
	openWithCustomInput
)

type IPCNavigateMsg struct {
	Path string
}

type Model struct {
	tree        paneTree
	preview     panePreview
	previewVP   viewport.Model
	focus       focusPane
	openWith    openWithState
	owCursor    int
	customInput string
	ipcActive   bool
	shellName   string
	width       int
	height      int
	quitting    bool
}

func New(startDir string, ipcActive bool) Model {
	if startDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			cwd = "/"
		}
		startDir = cwd
	}

	sh := shell.Detect()

	m := Model{
		tree:      newPaneTree(startDir),
		preview:   newPanePreview(),
		focus:     focusTree,
		ipcActive: ipcActive,
		shellName: sh.String(),
	}
	m.syncPreview()
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) CWD() string {
	return m.tree.cwd
}

func (m Model) SelectedName() string {
	e := m.tree.selected()
	if e == nil {
		return ""
	}
	return e.Name
}

func (m *Model) Navigate(path string) error {
	return m.tree.navigate(path)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.recalcLayout()

	case IPCNavigateMsg:
		m.tree.navigate(msg.Path)
		m.syncPreview()

	case tea.KeyMsg:
		if m.openWith != openWithOff {
			return m.updateOpenWith(msg)
		}
		return m.updateNormal(msg)
	}

	return m, nil
}

func (m Model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	action := mapKey(msg)

	switch action {
	case keyQuit:
		m.quitting = true
		return m, tea.Quit

	case keyUp:
		if m.focus == focusTree {
			m.tree.moveUp()
			m.syncPreview()
		} else {
			m.preview.scrollUp()
		}

	case keyDown:
		if m.focus == focusTree {
			m.tree.moveDown()
			m.syncPreview()
		} else {
			m.preview.scrollDown()
		}

	case keyTop:
		if m.focus == focusTree {
			m.tree.moveTop()
			m.syncPreview()
		}

	case keyBottom:
		if m.focus == focusTree {
			m.tree.moveBottom()
			m.syncPreview()
		}

	case keyLeft:
		m.tree.goUp()
		m.syncPreview()

	case keyRight:
		entered := m.tree.enterSelected()
		if !entered {
			m.focus = focusPreview
		}
		m.syncPreview()

	case keyOpenWith:
		e := m.tree.selected()
		if e != nil && e.Kind == clixfs.KindFile {
			m.openWith = openWithActive
			m.owCursor = 0
		}
	}

	return m, nil
}

func (m Model) updateOpenWith(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.openWith == openWithCustomInput {
		switch msg.String() {
		case "enter":
			return m.execOpen(opener.Custom, m.customInput)
		case "esc":
			m.openWith = openWithOff
			m.customInput = ""
		case "backspace":
			if len(m.customInput) > 0 {
				m.customInput = m.customInput[:len(m.customInput)-1]
			}
		default:
			if len(msg.Runes) == 1 {
				m.customInput += string(msg.Runes)
			}
		}
		return m, nil
	}

	options := owOptions()
	switch msg.String() {
	case "esc", "q":
		m.openWith = openWithOff
	case "h", "left":
		if m.owCursor > 0 {
			m.owCursor--
		}
	case "l", "right":
		if m.owCursor < len(options)-1 {
			m.owCursor++
		}
	case "enter":
		if m.owCursor == len(options)-1 {
			m.openWith = openWithCustomInput
			m.customInput = ""
			return m, nil
		}
		return m.execOpen(opener.Editor(m.owCursor), "")
	}
	return m, nil
}

func (m Model) execOpen(ed opener.Editor, custom string) (tea.Model, tea.Cmd) {
	e := m.tree.selected()
	if e == nil {
		m.openWith = openWithOff
		return m, nil
	}

	cmd := opener.Build(ed, custom, e.Path)
	m.openWith = openWithOff
	m.customInput = ""
	m.quitting = false

	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		return nil
	})
}

func (m *Model) syncPreview() {
	e := m.tree.selected()
	if e == nil || e.Kind == clixfs.KindDir {
		m.preview.clear()
		return
	}
	m.preview.load(e.Path)
}

func (m *Model) recalcLayout() {
	statusH := 1
	headerH := 0
	available := m.height - statusH - headerH

	treeW := m.width / 3
	if treeW < 20 {
		treeW = 20
	}
	previewW := m.width - treeW

	m.tree.width = treeW
	m.tree.height = available

	m.preview.width = previewW
	m.preview.height = available
}

func (m Model) View() string {
	if m.quitting {
		return ""
	}

	treePane := m.tree.render(m.focus == focusTree)
	previewPane := m.preview.render(m.focus == focusPreview)

	main := lipgloss.JoinHorizontal(lipgloss.Top, treePane, previewPane)

	status := m.renderStatus()

	return lipgloss.JoinVertical(lipgloss.Left, main, status)
}

func (m Model) renderStatus() string {
	w := m.width
	if w <= 0 {
		w = 80
	}

	var parts []string

	if m.openWith == openWithActive {
		parts = append(parts, m.renderOpenWithBar())
	} else if m.openWith == openWithCustomInput {
		parts = append(parts, fmt.Sprintf("open with custom: %s_", m.customInput))
	} else {
		parts = append(parts,
			styleDim.Render("[j/k] nav"),
			styleDim.Render("[h] up"),
			styleDim.Render("[l/enter] open"),
			styleDim.Render("[o] open with"),
			styleDim.Render("[q] quit"),
		)
	}

	ipcStr := styleDim.Render("ipc: off")
	if m.ipcActive {
		ipcStr = lipgloss.NewStyle().Foreground(colorIPCActive).Render("ipc: active")
	}

	shellStr := styleDim.Render("shell:" + m.shellName)

	left := strings.Join(parts, "  ")
	right := shellStr + "  " + ipcStr

	pad := w - lipgloss.Width(left) - lipgloss.Width(right)
	if pad < 1 {
		pad = 1
	}

	return left + strings.Repeat(" ", pad) + right
}

func (m Model) renderOpenWithBar() string {
	options := owOptions()
	var parts []string
	for i, o := range options {
		if i == m.owCursor {
			parts = append(parts, styleSelectedItem.Render(" "+o+" "))
		} else {
			parts = append(parts, styleDim.Render(" "+o+" "))
		}
	}
	return "open with: " + strings.Join(parts, " ")
}

func owOptions() []string {
	base := []string{"nano", "vim", "nvim", "cat", "custom"}
	return base
}

type execCmd struct {
	cmd *exec.Cmd
}

func (e execCmd) Run() error {
	e.cmd.Stdin = os.Stdin
	e.cmd.Stdout = os.Stdout
	e.cmd.Stderr = os.Stderr
	return e.cmd.Run()
}
