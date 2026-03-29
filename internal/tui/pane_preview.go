package tui

import (
	"fmt"
	"strings"

	clixfs "github.com/Lucas-Henry/clix/internal/fs"
)

type panePreview struct {
	preview  clixfs.Preview
	filePath string
	scroll   int
	width    int
	height   int
}

func newPanePreview() panePreview {
	return panePreview{width: 60, height: 20}
}

func (p *panePreview) load(path string) {
	p.filePath = path
	p.preview = clixfs.ReadPreview(path)
	p.scroll = 0
}

func (p *panePreview) clear() {
	p.filePath = ""
	p.preview = clixfs.Preview{}
	p.scroll = 0
}

func (p *panePreview) scrollDown() {
	maxScroll := len(p.preview.Lines) - p.visibleLines()
	if p.scroll < maxScroll {
		p.scroll++
	}
}

func (p *panePreview) scrollUp() {
	if p.scroll > 0 {
		p.scroll--
	}
}

func (p *panePreview) visibleLines() int {
	h := p.height - 4
	if h < 1 {
		return 1
	}
	return h
}

func (p panePreview) render(focused bool) string {
	innerW := p.width - 2
	innerH := p.height - 2
	if innerW < 4 {
		innerW = 4
	}
	if innerH < 1 {
		innerH = 1
	}

	var titleStr string
	if p.filePath == "" {
		titleStr = styleDim.Render("no file selected")
	} else {
		titleStr = styleHeader.Render(truncate(p.filePath, innerW))
	}
	titleLine := titleStr + "\n" + strings.Repeat("─", innerW)

	bodyHeight := innerH - 2
	var body string

	switch {
	case p.filePath == "":
		body = renderEmpty(bodyHeight)

	case p.preview.Error != "":
		body = styleDim.Render(p.preview.Error) + strings.Repeat("\n", bodyHeight-1)

	case p.preview.IsBinary:
		body = styleBinaryWarn.Render("[binary file]") + strings.Repeat("\n", bodyHeight-1)

	default:
		body = p.renderLines(innerW, bodyHeight)
	}

	content := titleLine + "\n" + body

	style := panelBorder.Width(innerW).Height(innerH)
	if focused {
		style = panelBorderFocused.Width(innerW).Height(innerH)
	}
	return style.Render(content)
}

func (p panePreview) renderLines(width, height int) string {
	if len(p.preview.Lines) == 0 {
		return styleDim.Render("[empty file]") + strings.Repeat("\n", height-1)
	}

	total := len(p.preview.Lines)
	numWidth := len(fmt.Sprintf("%d", total))
	contentW := width - numWidth - 2
	if contentW < 1 {
		contentW = 1
	}

	start := p.scroll
	end := start + height
	if end > total {
		end = total
	}

	var sb strings.Builder
	rendered := 0
	for i := start; i < end; i++ {
		lineNum := styleLineNum.Render(fmt.Sprintf("%*d ", numWidth, i+1))
		code := truncate(p.preview.Lines[i], contentW)
		sb.WriteString(lineNum + code + "\n")
		rendered++
	}

	for rendered < height {
		sb.WriteByte('\n')
		rendered++
	}

	return sb.String()
}

func renderEmpty(height int) string {
	return strings.Repeat("\n", height)
}
