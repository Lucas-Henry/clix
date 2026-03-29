package search

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFuzzyScore(t *testing.T) {
	cases := []struct {
		query  string
		target string
		wantGT int
	}{
		{"main", "main.go", 100},
		{"mai", "main.go", 1},
		{"app", "application.go", 1},
		{"xyz", "main.go", 0},
		{"", "main.go", 0},
	}

	for _, c := range cases {
		score := fuzzyScore(c.query, c.target)
		if c.wantGT == 0 && score != 0 {
			t.Errorf("fuzzyScore(%q, %q) = %d, want 0", c.query, c.target, score)
		}
		if c.wantGT > 0 && score < c.wantGT {
			t.Errorf("fuzzyScore(%q, %q) = %d, want > %d", c.query, c.target, score, c.wantGT)
		}
	}
}

func TestFuzzyNameSearch(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(root, "readme.md"), []byte("# readme"), 0644)
	sub := filepath.Join(root, "cmd")
	os.Mkdir(sub, 0755)
	os.WriteFile(filepath.Join(sub, "main.go"), []byte("package main"), 0644)

	results := FuzzyName("main", root, 3)
	if len(results) == 0 {
		t.Fatal("expected results for query 'main'")
	}
	for _, r := range results {
		if r.Score == 0 {
			t.Errorf("result %q has score 0", r.Name)
		}
	}

	empty := FuzzyName("", root, 3)
	if len(empty) != 0 {
		t.Error("empty query should return no results")
	}
}

func TestContentSearch(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "a.go"), []byte("package main\n\nfunc Hello() {}\n"), 0644)
	os.WriteFile(filepath.Join(root, "b.go"), []byte("package main\n\nfunc World() {}\n"), 0644)
	os.WriteFile(filepath.Join(root, "bin.dat"), []byte{0x00, 0x01, 0x02}, 0644)

	results := ContentSearch("Hello", root, 2)
	if len(results) == 0 {
		t.Fatal("expected content match for 'Hello'")
	}
	if results[0].Line == 0 {
		t.Error("line number should be non-zero")
	}

	noBin := ContentSearch("Hello", root, 2)
	for _, r := range noBin {
		if filepath.Base(r.Path) == "bin.dat" {
			t.Error("binary file should be skipped in content search")
		}
	}
}

func TestSortByScore(t *testing.T) {
	results := []NameMatch{
		{Name: "a", Score: 10},
		{Name: "b", Score: 150},
		{Name: "c", Score: 50},
	}
	sortByScore(results)
	if results[0].Score < results[1].Score || results[1].Score < results[2].Score {
		t.Errorf("results not sorted descending: %v %v %v",
			results[0].Score, results[1].Score, results[2].Score)
	}
}
