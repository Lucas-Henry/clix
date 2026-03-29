package tui

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	clixfs "github.com/Lucas-Henry/clix/internal/fs"
	"github.com/Lucas-Henry/clix/internal/opener"
	"github.com/Lucas-Henry/clix/internal/search"
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

type searchMode int

const (
	searchOff searchMode = iota
	searchName
	searchContent
)

type IPCNavigateMsg struct {
	Path string
}

type IPCSelectMsg struct {
	Name string
}

type searchDoneMsg struct {
	nameResults    []search.NameMatch
	contentResults []search.ContentMatch
}

type clickTimer struct {
	row  int
	time time.Time
}

const doubleClickThreshold = 400 * time.Millisecond

type Model struct {
	tree           paneTree
	preview        panePreview
	focus          focusPane
	openWith       openWithState
	owCursor       int
	customInput    string
	searchMode     searchMode
	searchInput    string
	nameResults    []search.NameMatch
	contentResults []search.ContentMatch
	resultCursor   int
	showResults    bool
	lastClick      clickTimer
	ipcActive      bool
	shellName      string
	width          int
	height         int
	quitting       bool
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

func (m *Model) SelectByName(name string) {
	for i, e := range m.tree.entries {
		if e.Name == name {
			m.tree.cursor = i
			m.syncPreview()
			return
		}
	}
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

	case IPCSelectMsg:
		m.SelectByName(msg.Name)

	case searchDoneMsg:
		m.nameResults = msg.nameResults
		m.contentResults = msg.contentResults
		m.resultCursor = 0
		m.showResults = true

	case tea.MouseMsg:
		return m.updateMouse(msg)

	case tea.KeyMsg:
		if m.searchMode != searchOff {
			return m.updateSearch(msg)
		}
		if m.openWith != openWithOff {
			return m.updateOpenWith(msg)
		}
		return m.updateNormal(msg)
	}

	return m, nil
}

func (m Model) updateMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		if m.focus == focusTree {
			m.tree.moveUp()
		} else {
			m.preview.scrollUp()
		}
		m.syncPreview()

	case tea.MouseButtonWheelDown:
		if m.focus == focusTree {
			m.tree.moveDown()
		} else {
			m.preview.scrollDown()
		}
		m.syncPreview()

	case tea.MouseButtonLeft:
		if msg.Action != tea.MouseActionRelease {
			break
		}
		treeW := m.tree.width
		if msg.X < treeW {
			// click in tree pane
			m.focus = focusTree
			idx := m.tree.itemAtRow(msg.Y)
			if idx < 0 {
				break
			}
			now := time.Now()
			isDouble := idx == m.lastClick.row && now.Sub(m.lastClick.time) < doubleClickThreshold
			m.lastClick = clickTimer{row: idx, time: now}
			if isDouble {
				m.tree.cursor = idx
				m.tree.enterSelected()
				m.syncPreview()
			} else {
				m.tree.cursor = idx
				m.syncPreview()
			}
		} else {
			// click in preview pane
			m.focus = focusPreview
		}
	}
	return m, nil
}

func (m Model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.showResults {
		return m.updateResults(msg)
	}

	action := mapKey(msg)

	switch action {
	case keyQuit:
		m.quitting = true
		return m, tea.Quit

	case keyEsc:
		m.showResults = false
		m.nameResults = nil
		m.contentResults = nil

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

	case keySearch:
		m.searchMode = searchName
		m.searchInput = ""
		m.showResults = false

	case keyContentSearch:
		m.searchMode = searchContent
		m.searchInput = ""
		m.showResults = false
	}

	return m, nil
}

func (m Model) updateSearch(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.searchMode = searchOff
		m.searchInput = ""
		m.showResults = false
	case "enter":
		mode := m.searchMode
		query := m.searchInput
		root := m.tree.cwd
		m.searchMode = searchOff
		return m, func() tea.Msg {
			var nm []search.NameMatch
			var cm []search.ContentMatch
			if mode == searchName {
				nm = search.FuzzyName(query, root, 5)
			} else {
				cm = search.ContentSearch(query, root, 5)
			}
			return searchDoneMsg{nameResults: nm, contentResults: cm}
		}
	case "backspace":
		if len(m.searchInput) > 0 {
			r := []rune(m.searchInput)
			m.searchInput = string(r[:len(r)-1])
		}
	default:
		if len(msg.Runes) > 0 {
			m.searchInput += string(msg.Runes)
		}
	}
	return m, nil
}

func (m Model) updateResults(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	total := len(m.nameResults) + len(m.contentResults)
	switch msg.String() {
	case "esc", "q":
		m.showResults = false
		m.nameResults = nil
		m.contentResults = nil
	case "j", "down":
		if m.resultCursor < total-1 {
			m.resultCursor++
		}
	case "k", "up":
		if m.resultCursor > 0 {
			m.resultCursor--
		}
	case "enter", "l":
		m.jumpToResult()
	}
	return m, nil
}

func (m *Model) jumpToResult() {
	if m.resultCursor < len(m.nameResults) {
		r := m.nameResults[m.resultCursor]
		if r.IsDir {
			m.tree.navigate(r.Path)
		} else {
			m.tree.navigate(r.Path[:len(r.Path)-len(r.Name)-1])
			m.SelectByName(r.Name)
		}
	} else {
		idx := m.resultCursor - len(m.nameResults)
		if idx < len(m.contentResults) {
			r := m.contentResults[idx]
			dir := r.Path[:len(r.Path)-len("/"+fmt.Sprintf("%s", r.Path[strings.LastIndex(r.Path, "/")+1:]))]
			name := r.Path[strings.LastIndex(r.Path, "/")+1:]
			_ = dir
			m.tree.navigate(r.Path[:strings.LastIndex(r.Path, "/")])
			m.SelectByName(name)
		}
	}
	m.showResults = false
	m.nameResults = nil
	m.contentResults = nil
	m.syncPreview()
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
	available := m.height - statusH

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

	if m.showResults {
		overlay := m.renderResults()
		return lipgloss.JoinVertical(lipgloss.Left, main, overlay, m.renderStatus())
	}

	return lipgloss.JoinVertical(lipgloss.Left, main, m.renderStatus())
}

func (m Model) renderResults() string {
	var sb strings.Builder
	sb.WriteString(styleDim.Render("results (enter to jump, esc to close):") + "\n")

	idx := 0
	for _, r := range m.nameResults {
		line := fmt.Sprintf("  %s", r.Path)
		if idx == m.resultCursor {
			sb.WriteString(styleSelectedItem.Render(line))
		} else {
			sb.WriteString(styleDim.Render(line))
		}
		sb.WriteByte('\n')
		idx++
	}
	for _, r := range m.contentResults {
		line := fmt.Sprintf("  %s:%d  %s", r.Path, r.Line, truncate(r.Content, 60))
		if idx == m.resultCursor {
			sb.WriteString(styleSelectedItem.Render(line))
		} else {
			sb.WriteString(styleDim.Render(line))
		}
		sb.WriteByte('\n')
		idx++
	}

	if idx == 0 {
		sb.WriteString(styleDim.Render("  no results found") + "\n")
	}

	return sb.String()
}

func (m Model) renderStatus() string {
	w := m.width
	if w <= 0 {
		w = 80
	}

	var left string

	switch {
	case m.searchMode == searchName:
		left = styleHeader.Render("/") + " " + m.searchInput + styleDim.Render("_")
	case m.searchMode == searchContent:
		left = styleHeader.Render("ctrl+f") + " " + m.searchInput + styleDim.Render("_")
	case m.openWith == openWithActive:
		left = m.renderOpenWithBar()
	case m.openWith == openWithCustomInput:
		left = styleHeader.Render("open: ") + m.customInput + styleDim.Render("_")
	default:
		parts := []string{
			styleDim.Render("[j/k] nav"),
			styleDim.Render("[h] up"),
			styleDim.Render("[l] open"),
			styleDim.Render("[o] open with"),
			styleDim.Render("[/] name search"),
			styleDim.Render("[ctrl+f] content"),
			styleDim.Render("[q] quit"),
		}
		left = strings.Join(parts, "  ")
	}

	ipcStr := styleDim.Render("ipc: off")
	if m.ipcActive {
		ipcStr = lipgloss.NewStyle().Foreground(colorIPCActive).Render("ipc: active")
	}
	shellStr := styleDim.Render("shell:" + m.shellName)
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
	return []string{"nano", "vim", "nvim", "cat", "custom"}
}
