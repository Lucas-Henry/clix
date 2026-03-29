package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	clixfs "github.com/Lucas-Henry/clix/internal/fs"
	"github.com/Lucas-Henry/clix/internal/icons"
	"github.com/charmbracelet/lipgloss"
)

type paneTree struct {
	cwd     string
	entries []clixfs.Entry
	cursor  int
	width   int
	height  int
	err     string
}

func newPaneTree(startDir string) paneTree {
	p := paneTree{width: 40, height: 20}
	p.load(startDir)
	return p
}

func (p *paneTree) load(dir string) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		p.err = err.Error()
		return
	}

	entries, err := clixfs.ReadDir(abs)
	if err != nil {
		p.err = err.Error()
		return
	}

	p.cwd = abs
	p.err = ""

	parent := clixfs.Entry{
		Name: "../",
		Path: filepath.Dir(abs),
		Kind: clixfs.KindDir,
	}
	p.entries = append([]clixfs.Entry{parent}, entries...)
	p.cursor = 0
}

func (p *paneTree) navigate(path string) error {
	p.load(path)
	if p.err != "" {
		return fmt.Errorf("%s", p.err)
	}
	return nil
}

func (p *paneTree) selected() *clixfs.Entry {
	if len(p.entries) == 0 {
		return nil
	}
	return &p.entries[p.cursor]
}

func (p *paneTree) moveUp() {
	if p.cursor > 0 {
		p.cursor--
	}
}

func (p *paneTree) moveDown() {
	if p.cursor < len(p.entries)-1 {
		p.cursor++
	}
}

func (p *paneTree) moveTop() {
	p.cursor = 0
}

func (p *paneTree) moveBottom() {
	if len(p.entries) > 0 {
		p.cursor = len(p.entries) - 1
	}
}

func (p *paneTree) enterSelected() bool {
	e := p.selected()
	if e == nil {
		return false
	}
	if e.Kind == clixfs.KindDir {
		p.load(e.Path)
		return true
	}
	return false
}

func (p *paneTree) goUp() {
	parent := filepath.Dir(p.cwd)
	if parent != p.cwd {
		prev := p.cwd
		p.load(parent)
		for i, e := range p.entries {
			if e.Path == prev {
				p.cursor = i
				break
			}
		}
	}
}

// itemRow returns the zero-based row index of a visible item given a y
// coordinate inside the panel (accounting for title and border).
func (p paneTree) itemAtRow(y int) int {
	// border(1) + title(1) + separator(1) = 3 lines before items
	itemY := y - 3
	if itemY < 0 {
		return -1
	}
	listHeight := p.height - 2 - 2 // inner minus title+sep
	start, end := visibleRange(p.cursor, len(p.entries), listHeight)
	idx := start + itemY
	if idx < start || idx >= end || idx >= len(p.entries) {
		return -1
	}
	return idx
}

func (p paneTree) render(focused bool) string {
	innerW := p.width - 2
	innerH := p.height - 2
	if innerW < 4 {
		innerW = 4
	}
	if innerH < 1 {
		innerH = 1
	}

	title := styleHeader.Render(truncate(p.cwd, innerW))
	titleLine := title + "\n" + strings.Repeat("─", innerW)

	listHeight := innerH - 2

	start, end := visibleRange(p.cursor, len(p.entries), listHeight)

	var sb strings.Builder
	for i := start; i < end && i < len(p.entries); i++ {
		e := p.entries[i]
		line := renderEntry(e, innerW-2)
		if i == p.cursor {
			line = styleSelectedItem.Width(innerW).Render(line)
		}
		sb.WriteString(line)
		sb.WriteByte('\n')
	}

	for rendered := end - start; rendered < listHeight; rendered++ {
		sb.WriteByte('\n')
	}

	content := titleLine + "\n" + sb.String()

	style := panelBorder.Width(innerW).Height(innerH)
	if focused {
		style = panelBorderFocused.Width(innerW).Height(innerH)
	}
	return style.Render(content)
}

func renderEntry(e clixfs.Entry, maxW int) string {
	icon := icons.ForEntry(e.Name, e.Kind == clixfs.KindDir, e.IsEmpty)

	var iconStyle lipgloss.Style
	var nameStr string

	switch e.Kind {
	case clixfs.KindDir:
		iconStyle = styleDir
		nameStr = styleDir.Render(truncate(e.Name, maxW-4)) + "/"
	case clixfs.KindSymlink:
		iconStyle = styleSymlink
		nameStr = styleSymlink.Render(truncate(e.Name, maxW-4))
	default:
		iconStyle = styleFile
		nameStr = styleFile.Render(truncate(e.Name, maxW-4))
	}

	return iconStyle.Render(icon) + " " + nameStr
}

func visibleRange(cursor, total, height int) (int, int) {
	if total == 0 || height <= 0 {
		return 0, 0
	}
	start := cursor - height/2
	if start < 0 {
		start = 0
	}
	end := start + height
	if end > total {
		end = total
		start = end - height
		if start < 0 {
			start = 0
		}
	}
	return start, end
}

func truncate(s string, max int) string {
	if max <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	if max <= 3 {
		return string(runes[:max])
	}
	return string(runes[:max-3]) + "..."
}
