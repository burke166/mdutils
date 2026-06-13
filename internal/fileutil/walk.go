package fileutil

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var ignoredDirNames = map[string]struct{}{
	".git":         {},
	"node_modules": {},
	"vendor":       {},
	"bin":          {},
	"obj":          {},
	"dist":         {},
	"build":        {},
}

// CollectMarkdownFiles returns Markdown files from a single file or directory.
// When path is a directory, recursive controls whether subdirectories are walked.
func CollectMarkdownFiles(path string, recursive bool) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		if !isMarkdownFile(path) {
			return nil, fmt.Errorf("not a Markdown file: %s", path)
		}
		return []string{path}, nil
	}

	if recursive {
		return collectMarkdownFilesRecursive(path)
	}

	return collectMarkdownFilesFlat(path)
}

func collectMarkdownFilesFlat(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		if isMarkdownFile(path) {
			files = append(files, path)
		}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no Markdown files found in directory: %s", dir)
	}

	sort.Strings(files)
	return files, nil
}

func collectMarkdownFilesRecursive(root string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path != root && shouldIgnoreDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if isMarkdownFile(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no Markdown files found in directory: %s", root)
	}

	sort.Strings(files)
	return files, nil
}

func shouldIgnoreDir(name string) bool {
	_, ok := ignoredDirNames[name]
	return ok
}

func isMarkdownFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".markdown"
}

// FilterExcluded removes files matching any of the glob patterns.
func FilterExcluded(files []string, root string, patterns []string) []string {
	if len(patterns) == 0 {
		return files
	}

	filtered := make([]string, 0, len(files))
	for _, file := range files {
		if isExcluded(file, root, patterns) {
			continue
		}
		filtered = append(filtered, file)
	}
	return filtered
}

func isExcluded(file, root string, patterns []string) bool {
	rel := file
	if root != "" {
		if relative, err := filepath.Rel(root, file); err == nil {
			rel = relative
		}
	}

	rel = filepath.ToSlash(rel)
	base := filepath.Base(file)

	for _, pattern := range patterns {
		pattern = filepath.ToSlash(pattern)
		if matchGlob(pattern, rel) || matchGlob(pattern, base) {
			return true
		}
	}

	return false
}

func matchGlob(pattern, name string) bool {
	if pattern == "" {
		return false
	}

	if strings.Contains(pattern, "**") {
		return matchDoubleStar(pattern, name)
	}

	matched, err := filepath.Match(pattern, name)
	return err == nil && matched
}

func matchDoubleStar(pattern, name string) bool {
	if pattern == "**" {
		return true
	}

	parts := strings.Split(pattern, "**")
	if len(parts) == 1 {
		matched, err := filepath.Match(pattern, name)
		return err == nil && matched
	}

	prefix := strings.TrimSuffix(parts[0], "/")
	suffix := strings.TrimPrefix(parts[1], "/")

	if prefix != "" {
		if name == prefix {
			return suffix == ""
		}
		if !strings.HasPrefix(name, prefix+"/") {
			return false
		}
		name = strings.TrimPrefix(name, prefix+"/")
	}

	if suffix == "" {
		return true
	}

	if strings.Contains(suffix, "/") {
		return strings.HasPrefix(name, suffix) || strings.HasSuffix(name, "/"+suffix)
	}

	matched, err := filepath.Match(suffix, name)
	if err == nil && matched {
		return true
	}

	for _, part := range strings.Split(name, "/") {
		if matched, err := filepath.Match(suffix, part); err == nil && matched {
			return true
		}
	}

	return false
}
