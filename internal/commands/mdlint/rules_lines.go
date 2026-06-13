package mdlint

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

func lintLines(file string, content string, cfg Config) []Issue {
	rules := cfg.Rules
	if !rules.NoTrailingWhitespace && rules.MaxLineLength <= 0 &&
		!rules.NoMultipleBlankLines && !rules.RequireCodeFenceLanguage {
		return nil
	}

	lines := markdown.SplitLines(content)
	var issues []Issue

	inFence := false
	var fenceChar byte
	prevBlank := false

	for lineIndex, line := range lines {
		lineNumber := lineIndex + 1
		body := strings.TrimRight(line, "\r\n")

		isFence, char := markdown.IsFenceLine(line)
		if isFence {
			if !inFence {
				if rules.RequireCodeFenceLanguage && !fenceHasLanguage(body) {
					issues = append(issues, Issue{
						File:    file,
						Line:    lineNumber,
						Column:  1,
						RuleID:  RuleRequireCodeFenceLang,
						Message: "fenced code block is missing a language identifier",
					})
				}
				inFence = true
				fenceChar = char
			} else if char == fenceChar {
				inFence = false
				fenceChar = 0
			}
			prevBlank = false
			continue
		}

		if inFence {
			prevBlank = false
			continue
		}

		if rules.NoTrailingWhitespace {
			if col := trailingWhitespaceColumn(body); col > 0 {
				issues = append(issues, Issue{
					File:    file,
					Line:    lineNumber,
					Column:  col,
					RuleID:  RuleNoTrailingWhitespace,
					Message: "line has trailing whitespace",
				})
			}
		}

		if rules.MaxLineLength > 0 {
			if col := lineLengthExceededColumn(body, rules.MaxLineLength); col > 0 {
				issues = append(issues, Issue{
					File:    file,
					Line:    lineNumber,
					Column:  col,
					RuleID:  RuleMaxLineLength,
					Message: fmt.Sprintf("line exceeds maximum length of %d characters", rules.MaxLineLength),
				})
			}
		}

		blank := markdown.IsBlankLine(line)
		if rules.NoMultipleBlankLines && blank && prevBlank {
			issues = append(issues, Issue{
				File:    file,
				Line:    lineNumber,
				Column:  1,
				RuleID:  RuleNoMultipleBlankLines,
				Message: "more than one consecutive blank line",
			})
		}
		prevBlank = blank
	}

	for i := range issues {
		issues[i] = withSeverity(cfg, issues[i])
	}

	return issues
}

func fenceHasLanguage(line string) bool {
	trimmed := strings.TrimLeft(line, " \t")
	trimmed = strings.TrimRight(trimmed, "\r\n")

	switch {
	case strings.HasPrefix(trimmed, "```"):
		rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
		return rest != ""
	case strings.HasPrefix(trimmed, "~~~"):
		rest := strings.TrimSpace(strings.TrimPrefix(trimmed, "~~~"))
		return rest != ""
	default:
		return true
	}
}

func trailingWhitespaceColumn(line string) int {
	for i := len(line) - 1; i >= 0; i-- {
		if line[i] != ' ' && line[i] != '\t' {
			if i == len(line)-1 {
				return 0
			}
			return i + 2
		}
	}
	return 0
}

func lineLengthExceededColumn(line string, max int) int {
	if utf8.RuneCountInString(line) <= max {
		return 0
	}

	count := 0
	for i := range line {
		count++
		if count > max {
			return i + 1
		}
	}
	return 0
}
