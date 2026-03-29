package fs

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type EntryKind int

const (
	KindDir EntryKind = iota
	KindFile
	KindSymlink
)

type Entry struct {
	Name    string
	Path    string
	Kind    EntryKind
	Size    int64
	IsEmpty bool
}

func ReadDir(dir string) ([]Entry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	dirs := make([]Entry, 0, len(entries))
	files := make([]Entry, 0, len(entries))

	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}

		kind := KindFile
		if e.IsDir() {
			kind = KindDir
		} else if info.Mode()&os.ModeSymlink != 0 {
			kind = KindSymlink
		}

		entry := Entry{
			Name: e.Name(),
			Path: filepath.Join(dir, e.Name()),
			Kind: kind,
			Size: info.Size(),
		}

		if kind == KindDir {
			entry.IsEmpty = isDirEmpty(entry.Path)
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}

	sort.Slice(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i].Name) < strings.ToLower(dirs[j].Name)
	})
	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	return append(dirs, files...), nil
}

func isDirEmpty(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	return len(entries) == 0
}
