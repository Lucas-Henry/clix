package fs

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	maxPreviewLines = 200
	sniffBytes      = 512
)

type Preview struct {
	Lines    []string
	IsBinary bool
	Ext      string
	Error    string
}

func ReadPreview(path string) Preview {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(path), "."))
	p := Preview{Ext: ext}

	f, err := os.Open(path)
	if err != nil {
		p.Error = fmt.Sprintf("cannot open: %v", err)
		return p
	}
	defer f.Close()

	sniff := make([]byte, sniffBytes)
	n, _ := f.Read(sniff)
	if isBinary(sniff[:n]) {
		p.IsBinary = true
		return p
	}

	if _, err := f.Seek(0, 0); err != nil {
		p.Error = fmt.Sprintf("seek error: %v", err)
		return p
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() && len(p.Lines) < maxPreviewLines {
		p.Lines = append(p.Lines, scanner.Text())
	}

	return p
}

func isBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	return bytes.IndexByte(data, 0) != -1
}
