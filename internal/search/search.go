package search

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type NameMatch struct {
	Path  string
	Name  string
	IsDir bool
	Score int
}

type ContentMatch struct {
	Path    string
	Line    int
	Content string
}

func FuzzyName(query, root string, maxDepth int) []NameMatch {
	if query == "" {
		return nil
	}
	lower := strings.ToLower(query)
	var results []NameMatch
	walkDepth(root, 0, maxDepth, func(path string, isDir bool) bool {
		name := filepath.Base(path)
		score := fuzzyScore(lower, strings.ToLower(name))
		if score > 0 {
			results = append(results, NameMatch{
				Path:  path,
				Name:  name,
				IsDir: isDir,
				Score: score,
			})
		}
		return true
	})
	sortByScore(results)
	return results
}

func ContentSearch(query, root string, maxDepth int) []ContentMatch {
	if query == "" {
		return nil
	}
	lower := strings.ToLower(query)
	var results []ContentMatch
	walkDepth(root, 0, maxDepth, func(path string, isDir bool) bool {
		if isDir {
			return true
		}
		matches := grepFile(path, lower)
		results = append(results, matches...)
		return true
	})
	return results
}

func grepFile(path, lowerQuery string) []ContentMatch {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	sniff := make([]byte, 512)
	n, _ := f.Read(sniff)
	for _, b := range sniff[:n] {
		if b == 0 {
			return nil
		}
	}
	if _, err := f.Seek(0, 0); err != nil {
		return nil
	}

	var matches []ContentMatch
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), lowerQuery) {
			matches = append(matches, ContentMatch{
				Path:    path,
				Line:    lineNum,
				Content: strings.TrimSpace(line),
			})
			if len(matches) >= 5 {
				break
			}
		}
	}
	return matches
}

func walkDepth(dir string, depth, maxDepth int, fn func(path string, isDir bool) bool) {
	if depth > maxDepth {
		return
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		path := filepath.Join(dir, e.Name())
		if !fn(path, e.IsDir()) {
			return
		}
		if e.IsDir() {
			walkDepth(path, depth+1, maxDepth, fn)
		}
	}
}

func fuzzyScore(query, target string) int {
	if query == "" {
		return 0
	}

	if strings.Contains(target, query) {
		score := 100
		if target == query {
			score += 100
		}
		if strings.HasPrefix(target, query) {
			score += 50
		}
		return score
	}

	qi := 0
	qRunes := []rune(query)
	tRunes := []rune(target)
	lastMatch := -1
	consecutive := 0
	score := 0

	for ti, tr := range tRunes {
		if qi >= len(qRunes) {
			break
		}
		if unicode.ToLower(tr) == unicode.ToLower(qRunes[qi]) {
			if lastMatch == ti-1 {
				consecutive++
				score += 10 + consecutive*5
			} else {
				consecutive = 0
				score += 5
			}
			if ti == 0 || tRunes[ti-1] == '_' || tRunes[ti-1] == '-' || tRunes[ti-1] == '.' {
				score += 15
			}
			lastMatch = ti
			qi++
		}
	}

	if qi == len(qRunes) {
		return score
	}
	return 0
}

func sortByScore(results []NameMatch) {
	n := len(results)
	for i := 1; i < n; i++ {
		key := results[i]
		j := i - 1
		for j >= 0 && results[j].Score < key.Score {
			results[j+1] = results[j]
			j--
		}
		results[j+1] = key
	}
}
