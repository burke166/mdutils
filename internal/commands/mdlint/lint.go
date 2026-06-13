package mdlint

import (
	"os"
	"sort"
)

func LintFile(path string, cfg Config) ([]Issue, error) {
	source, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return LintContent(path, source, cfg)
}

func LintContent(file string, source []byte, cfg Config) ([]Issue, error) {
	content := string(source)

	headingIssues, err := lintHeadings(file, source, cfg)
	if err != nil {
		return nil, err
	}

	issues := append(headingIssues, lintLines(file, content, cfg)...)
	issues = append(issues, lintLinks(file, content, cfg)...)

	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Line != issues[j].Line {
			return issues[i].Line < issues[j].Line
		}
		if issues[i].Column != issues[j].Column {
			return issues[i].Column < issues[j].Column
		}
		return issues[i].RuleID < issues[j].RuleID
	})

	return issues, nil
}

func LintPaths(paths []string, cfg Config) ([]Issue, error) {
	var issues []Issue

	for _, path := range paths {
		fileIssues, err := LintFile(path, cfg)
		if err != nil {
			return nil, err
		}
		issues = append(issues, fileIssues...)
	}

	return issues, nil
}
