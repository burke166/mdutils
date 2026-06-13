package mdlint

import (
	"regexp"
	"strings"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

var (
	linkPattern  = regexp.MustCompile(`\[([^\]]*)\]\(([^)]*)\)`)
	imagePattern = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]*)\)`)
)

func lintLinks(file string, content string, cfg Config) []Issue {
	rules := cfg.Rules
	if !rules.NoEmptyLinks && !rules.RequireImageAltText {
		return nil
	}

	lines := markdown.SplitLines(content)
	var issues []Issue

	inFence := false
	var fenceChar byte

	for lineIndex, line := range lines {
		lineNumber := lineIndex + 1

		isFence, char := markdown.IsFenceLine(line)
		if isFence {
			if !inFence {
				inFence = true
				fenceChar = char
			} else if char == fenceChar {
				inFence = false
				fenceChar = 0
			}
			continue
		}
		if inFence {
			continue
		}

		body := strings.TrimRight(line, "\r\n")

		if rules.NoEmptyLinks {
			for _, match := range linkPattern.FindAllStringSubmatchIndex(body, -1) {
				url := body[match[4]:match[5]]
				if strings.TrimSpace(url) != "" {
					continue
				}
				issues = append(issues, Issue{
					File:    file,
					Line:    lineNumber,
					Column:  match[0] + 1,
					RuleID:  RuleNoEmptyLinks,
					Message: "link target is empty",
				})
			}
		}

		if rules.RequireImageAltText {
			for _, match := range imagePattern.FindAllStringSubmatchIndex(body, -1) {
				alt := body[match[2]:match[3]]
				if strings.TrimSpace(alt) != "" {
					continue
				}
				issues = append(issues, Issue{
					File:    file,
					Line:    lineNumber,
					Column:  match[0] + 1,
					RuleID:  RuleRequireImageAltText,
					Message: "image is missing alt text",
				})
			}
		}
	}

	for i := range issues {
		issues[i] = withSeverity(cfg, issues[i])
	}

	return issues
}
