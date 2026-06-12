package mdmerge

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func MergeContents(contents []string) string {
	if len(contents) == 0 {
		return ""
	}

	trimmed := make([]string, len(contents))
	for i, content := range contents {
		trimmed[i] = strings.Trim(content, "\n")
	}

	return strings.Join(trimmed, "\n\n") + "\n"
}

func CollectMarkdownFiles(args []string) ([]string, error) {
	var files []string

	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			dirFiles, err := listMarkdownFiles(arg)
			if err != nil {
				return nil, err
			}
			if len(dirFiles) == 0 {
				return nil, fmt.Errorf("no Markdown files found in directory: %s", arg)
			}
			files = append(files, dirFiles...)
			continue
		}

		if !strings.EqualFold(filepath.Ext(arg), ".md") {
			return nil, fmt.Errorf("not a Markdown file: %s", arg)
		}

		files = append(files, arg)
	}

	if len(files) == 0 {
		return nil, errors.New("no input files specified")
	}

	return files, nil
}

func listMarkdownFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.EqualFold(filepath.Ext(entry.Name()), ".md") {
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}

	sort.Strings(files)
	return files, nil
}

func ReadMarkdownFiles(paths []string) ([]string, error) {
	contents := make([]string, 0, len(paths))
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		contents = append(contents, string(data))
	}
	return contents, nil
}

func IsNotExist(err error) bool {
	return errors.Is(err, fs.ErrNotExist)
}
