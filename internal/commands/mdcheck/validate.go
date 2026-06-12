package mdcheck

import (
	"fmt"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

const (
	RuleMissingH1        = "missing-h1"
	RuleMultipleH1       = "multiple-h1"
	RuleSkippedLevel     = "skipped-level"
	RuleDuplicateHeading = "duplicate-heading"
	RuleEmptyHeading     = "empty-heading"
	RuleMaxLevel         = "max-level"
)

type Diagnostic struct {
	Rule    string `json:"rule"`
	Message string `json:"message"`
	Line    int    `json:"line"`
	Level   int    `json:"level"`
	Text    string `json:"text"`
}

type ValidateOptions struct {
	CheckDuplicates bool
	AllowMultipleH1 bool
	MaxLevel        int
}

func Validate(headings []markdown.Heading, opts ValidateOptions) []Diagnostic {
	var diags []Diagnostic

	diags = append(diags, checkMissingH1(headings)...)
	diags = append(diags, checkMultipleH1(headings, opts.AllowMultipleH1)...)
	diags = append(diags, checkSkippedLevels(headings)...)
	if opts.CheckDuplicates {
		diags = append(diags, checkDuplicateHeadings(headings)...)
	}
	diags = append(diags, checkEmptyHeadings(headings)...)
	if opts.MaxLevel > 0 {
		diags = append(diags, checkMaxLevel(headings, opts.MaxLevel)...)
	}

	return diags
}

func checkMissingH1(headings []markdown.Heading) []Diagnostic {
	for _, h := range headings {
		if h.Level == 1 {
			return nil
		}
	}

	return []Diagnostic{{
		Rule:    RuleMissingH1,
		Message: "document has no H1 heading",
	}}
}

func checkMultipleH1(headings []markdown.Heading, allowMultiple bool) []Diagnostic {
	if allowMultiple {
		return nil
	}

	var diags []Diagnostic
	var first bool

	for _, h := range headings {
		if h.Level != 1 {
			continue
		}

		if !first {
			first = true
			continue
		}

		diags = append(diags, Diagnostic{
			Rule:    RuleMultipleH1,
			Message: "document has more than one H1 heading",
			Line:    h.Line,
			Level:   h.Level,
			Text:    h.Text,
		})
	}

	return diags
}

func checkSkippedLevels(headings []markdown.Heading) []Diagnostic {
	var diags []Diagnostic
	prevLevel := 0

	for _, h := range headings {
		if prevLevel > 0 && h.Level > prevLevel+1 {
			diags = append(diags, Diagnostic{
				Rule:    RuleSkippedLevel,
				Message: fmt.Sprintf("skipped heading level (H%d followed by H%d)", prevLevel, h.Level),
				Line:    h.Line,
				Level:   h.Level,
				Text:    h.Text,
			})
		}

		prevLevel = h.Level
	}

	return diags
}

func checkDuplicateHeadings(headings []markdown.Heading) []Diagnostic {
	seen := make(map[string]bool)
	var diags []Diagnostic

	for _, h := range headings {
		if h.Text == "" {
			continue
		}

		if seen[h.Text] {
			diags = append(diags, Diagnostic{
				Rule:    RuleDuplicateHeading,
				Message: fmt.Sprintf("duplicate heading text %q", h.Text),
				Line:    h.Line,
				Level:   h.Level,
				Text:    h.Text,
			})
			continue
		}

		seen[h.Text] = true
	}

	return diags
}

func checkEmptyHeadings(headings []markdown.Heading) []Diagnostic {
	var diags []Diagnostic

	for _, h := range headings {
		if h.Text != "" {
			continue
		}

		diags = append(diags, Diagnostic{
			Rule:    RuleEmptyHeading,
			Message: "empty heading",
			Line:    h.Line,
			Level:   h.Level,
		})
	}

	return diags
}

func checkMaxLevel(headings []markdown.Heading, maxLevel int) []Diagnostic {
	var diags []Diagnostic

	for _, h := range headings {
		if h.Level <= maxLevel {
			continue
		}

		diags = append(diags, Diagnostic{
			Rule:    RuleMaxLevel,
			Message: fmt.Sprintf("heading level H%d exceeds maximum allowed H%d", h.Level, maxLevel),
			Line:    h.Line,
			Level:   h.Level,
			Text:    h.Text,
		})
	}

	return diags
}
