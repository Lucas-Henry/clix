package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	tmp := t.TempDir()

	os.Mkdir(filepath.Join(tmp, "alpha"), 0755)
	os.Mkdir(filepath.Join(tmp, "beta"), 0755)
	os.WriteFile(filepath.Join(tmp, "file.txt"), []byte("hello"), 0644)
	os.WriteFile(filepath.Join(tmp, "readme.md"), []byte("world"), 0644)

	entries, err := ReadDir(tmp)
	if err != nil {
		t.Fatalf("ReadDir error: %v", err)
	}

	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}

	if entries[0].Kind != KindDir || entries[1].Kind != KindDir {
		t.Error("expected directories to come first")
	}

	if entries[0].Name > entries[1].Name {
		t.Errorf("dirs not sorted: %q > %q", entries[0].Name, entries[1].Name)
	}

	if entries[2].Kind != KindFile || entries[3].Kind != KindFile {
		t.Error("expected files after directories")
	}
}

func TestReadPreviewText(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test.go")
	content := "package main\n\nfunc main() {}\n"
	os.WriteFile(path, []byte(content), 0644)

	p := ReadPreview(path)
	if p.IsBinary {
		t.Error("expected text file, got binary")
	}
	if p.Error != "" {
		t.Errorf("unexpected error: %s", p.Error)
	}
	if len(p.Lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(p.Lines))
	}
}

func TestReadPreviewBinary(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "bin.dat")
	os.WriteFile(path, []byte{0x00, 0x01, 0x02, 0x00, 0xFF}, 0644)

	p := ReadPreview(path)
	if !p.IsBinary {
		t.Error("expected binary detection")
	}
}

func TestReadPreviewMissing(t *testing.T) {
	p := ReadPreview("/nonexistent/path/file.txt")
	if p.Error == "" {
		t.Error("expected error for missing file")
	}
}

func TestIsDirEmpty(t *testing.T) {
	tmp := t.TempDir()

	if !isDirEmpty(tmp) {
		t.Error("new temp dir should be empty")
	}

	os.WriteFile(filepath.Join(tmp, "x"), []byte("x"), 0644)
	if isDirEmpty(tmp) {
		t.Error("dir with file should not be empty")
	}
}
