package mdlint

import (
	"fmt"
	"github.com/computercodeblue/mdutils/internal/markdown"
)

func lintHeadings(file string, source []byte, cfg Config) ([]Issue, error) {
	rules := cfg.Rules
	if !rules.SingleH1 && !rules.NoMissingH1 && !rules.NoSkippedHeadingLevels &&
		!rules.NoDuplicateHeadings && !rules.NoEmptyHeadings && !rules.NoEmptySections &&
		rules.MaxHeadingLength <= 0 {
		return nil, nil
	}

	headings, err := markdown.ExtractHeadings(source)
	if err != nil {
		return nil, err
	}

	var issues []Issue

	if rules.NoMissingH1 {
		issues = append(issues, checkMissingH1(file, headings)...)
	}
	if rules.SingleH1 {
		issues = append(issues, checkMultipleH1(file, headings)...)
	}
	if rules.NoSkippedHeadingLevels {
		issues = append(issues, checkSkippedHeadingLevels(file, headings)...)
	}
	if rules.NoDuplicateHeadings {
		issues = append(issues, checkDuplicateHeadings(file, headings)...)
	}
	if rules.NoEmptyHeadings {
		issues = append(issues, checkEmptyHeadings(file, headings)...)
	}
	if rules.MaxHeadingLength > 0 {
		issues = append(issues, checkMaxHeadingLength(file, headings, rules.MaxHeadingLength)...)
	}
	if rules.NoEmptySections {
		issues = append(issues, checkEmptySections(file, source, headings)...)
	}

	for i := range issues {
		issues[i] = withSeverity(cfg, issues[i])
	}

	return issues, nil
}

func checkMissingH1(file string, headings []markdown.Heading) []Issue {
	for _, h := range headings {
		if h.Level == 1 {
			return nil
		}
	}

	return []Issue{{
		File:    file,
		Line:    1,
		RuleID:  RuleNoMissingH1,
		Message: "document has no H1 heading",
	}}
}

func checkMultipleH1(file string, headings []markdown.Heading) []Issue {
	var issues []Issue
	var seen bool

	for _, h := range headings {
		if h.Level != 1 {
			continue
		}
		if !seen {
			seen = true
			continue
		}

		issues = append(issues, Issue{
			File:    file,
			Line:    h.Line,
			RuleID:  RuleSingleH1,
			Message: "document has multiple H1 headings",
		})
	}

	return issues
}

func checkSkippedHeadingLevels(file string, headings []markdown.Heading) []Issue {
	var issues []Issue
	prevLevel := 0

	for _, h := range headings {
		if prevLevel > 0 && h.Level > prevLevel+1 {
			issues = append(issues, Issue{
				File:   file,
				Line:   h.Line,
				RuleID: RuleNoSkippedHeadingLevels,
				Message: fmt.Sprintf(
					"heading level skipped from H%d to H%d",
					prevLevel,
					h.Level,
				),
			})
		}
		prevLevel = h.Level
	}

	return issues
}

func checkDuplicateHeadings(file string, headings []markdown.Heading) []Issue {
	seen := make(map[string]bool)
	var issues []Issue

	for _, h := range headings {
		if h.Text == "" {
			continue
		}
		if seen[h.Text] {
			issues = append(issues, Issue{
				File:    file,
				Line:    h.Line,
				RuleID:  RuleNoDuplicateHeadings,
				Message: fmt.Sprintf("duplicate heading %q", h.Text),
			})
			continue
		}
		seen[h.Text] = true
	}

	return issues
}

func checkEmptyHeadings(file string, headings []markdown.Heading) []Issue {
	var issues []Issue

	for _, h := range headings {
		if h.Text != "" {
			continue
		}
		issues = append(issues, Issue{
			File:    file,
			Line:    h.Line,
			RuleID:  RuleNoEmptyHeadings,
			Message: "heading has no text",
		})
	}

	return issues
}

func checkMaxHeadingLength(file string, headings []markdown.Heading, max int) []Issue {
	var issues []Issue

	for _, h := range headings {
		if len(h.Text) <= max {
			continue
		}
		issues = append(issues, Issue{
			File:    file,
			Line:    h.Line,
			Column:  max + 1,
			RuleID:  RuleMaxHeadingLength,
			Message: fmt.Sprintf("heading exceeds maximum length of %d characters", max),
		})
	}

	return issues
}

func checkEmptySections(file string, source []byte, headings []markdown.Heading) []Issue {
	if len(headings) == 0 {
		return nil
	}

	lines := markdown.SplitLines(string(source))
	var issues []Issue

	for i, h := range headings {
		start := h.Line
		if start < 1 || start > len(lines) {
			continue
		}

		end := len(lines)
		if i+1 < len(headings) {
			end = headings[i+1].Line - 1
		}

		if sectionIsEmpty(lines, start, end) {
			issues = append(issues, Issue{
				File:    file,
				Line:    h.Line,
				RuleID:  RuleNoEmptySections,
				Message: "section has no content",
			})
		}
	}

	return issues
}

func sectionIsEmpty(lines []string, headingLine, endLine int) bool {
	for i := headingLine; i < endLine; i++ {
		if i < 0 || i >= len(lines) {
			continue
		}
		if !markdown.IsBlankLine(lines[i]) {
			return false
		}
	}
	return true
}
