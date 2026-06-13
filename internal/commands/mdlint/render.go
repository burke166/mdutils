package mdlint

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func RenderText(issues []Issue, quiet bool) string {
	if len(issues) == 0 {
		return ""
	}

	display := filterIssues(issues, quiet)
	if len(display) == 0 {
		return ""
	}

	byFile := groupIssuesByFile(display)
	files := make([]string, 0, len(byFile))
	for file := range byFile {
		files = append(files, file)
	}
	sort.Strings(files)

	var b strings.Builder
	for fileIndex, file := range files {
		if fileIndex > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(file)
		b.WriteByte('\n')

		for _, issue := range byFile[file] {
			b.WriteString("  ")
			b.WriteString(formatIssueLine(issue))
			b.WriteByte('\n')
		}
	}

	b.WriteByte('\n')
	b.WriteString(fmt.Sprintf("%d issues found\n", len(issues)))

	return b.String()
}

func RenderJSON(issues []Issue, quiet bool) (string, error) {
	display := filterIssues(issues, quiet)
	if display == nil {
		display = []Issue{}
	}

	data, err := json.MarshalIndent(display, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data) + "\n", nil
}

func filterIssues(issues []Issue, quiet bool) []Issue {
	if !quiet {
		return issues
	}

	filtered := make([]Issue, 0, len(issues))
	for _, issue := range issues {
		if issue.Severity == SeverityError {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

func groupIssuesByFile(issues []Issue) map[string][]Issue {
	byFile := make(map[string][]Issue)
	for _, issue := range issues {
		byFile[issue.File] = append(byFile[issue.File], issue)
	}
	return byFile
}

func formatIssueLine(issue Issue) string {
	location := fmt.Sprintf("line %d", issue.Line)
	if issue.Column > 0 {
		location = fmt.Sprintf("line %d:%d", issue.Line, issue.Column)
	}

	return fmt.Sprintf("%s: %s %s: %s", location, issue.Severity, issue.RuleID, issue.Message)
}
