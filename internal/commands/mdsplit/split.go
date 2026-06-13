package mdsplit

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

type Section struct {
	Heading string
	Slug    string
	Content string
}

func SplitMarkdown(content string, level int) ([]Section, error) {
	if level < 1 || level > 6 {
		return nil, fmt.Errorf("heading level must be between 1 and 6")
	}

	lines := markdown.SplitLines(content)
	var sections []Section
	var current *Section
	var preamble []string

	inFence := false
	var fenceChar byte
	found := false

	flushPreamble := func() {
		if len(preamble) == 0 {
			return
		}
		sections = append(sections, Section{
			Content: joinLines(preamble),
		})
		preamble = nil
	}

	startSection := func(heading, line string) {
		flushPreamble()
		if current != nil {
			sections = append(sections, *current)
		}
		current = &Section{
			Heading: heading,
			Slug:    SlugifyFilename(heading),
		}
		current.Content = line
		found = true
	}

	for _, line := range lines {
		isFence, char := markdown.IsFenceLine(line)
		if isFence {
			if !inFence {
				inFence = true
				fenceChar = char
			} else if char == fenceChar {
				inFence = false
				fenceChar = 0
			}
		} else if !inFence {
			if headingLevel(line) == level {
				heading := atxHeadingText(line)
				startSection(heading, line)
				continue
			}
		}

		if current != nil {
			current.Content = appendLine(current.Content, line)
		} else {
			preamble = append(preamble, line)
		}
	}

	if current != nil {
		sections = append(sections, *current)
	} else {
		flushPreamble()
	}

	for i := range sections {
		sections[i].Content = trimSectionContent(sections[i].Content)
	}

	if !found {
		return nil, errors.New("no matching headings found")
	}

	return sections, nil
}

func joinLines(lines []string) string {
	return strings.Join(lines, "")
}

func appendLine(content, line string) string {
	return content + line
}

func trimSectionContent(content string) string {
	if content == "" {
		return content
	}
	return strings.TrimRight(content, "\n") + "\n"
}

func headingLevel(line string) int {
	trimmed := strings.TrimLeft(line, " \t")
	if len(trimmed) == 0 || trimmed[0] != '#' {
		return 0
	}

	hashes := 0
	for hashes < len(trimmed) && trimmed[hashes] == '#' {
		hashes++
	}
	if hashes > 6 {
		return 0
	}

	rest := trimmed[hashes:]
	if rest == "" {
		return hashes
	}
	if rest[0] == ' ' || rest[0] == '\t' {
		return hashes
	}

	return 0
}

func atxHeadingText(line string) string {
	trimmed := strings.TrimLeft(line, " \t")
	i := 0
	for i < len(trimmed) && trimmed[i] == '#' {
		i++
	}
	text := strings.TrimSpace(trimmed[i:])
	return strings.TrimRight(text, " #")
}

func SlugifyFilename(heading string) string {
	s := strings.ToLower(strings.TrimSpace(heading))

	var b strings.Builder
	lastWasDash := false

	for _, r := range s {
		switch {
		case isIllegalFilenameRune(r):
			if b.Len() > 0 && !lastWasDash {
				b.WriteByte('-')
				lastWasDash = true
			}
		case unicode.IsSpace(r):
			if b.Len() > 0 && !lastWasDash {
				b.WriteByte('-')
				lastWasDash = true
			}
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastWasDash = false
		default:
			if b.Len() > 0 && !lastWasDash {
				b.WriteByte('-')
				lastWasDash = true
			}
		}
	}

	slug := strings.Trim(b.String(), "-")
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	if slug == "" {
		return "section"
	}
	return slug
}

func isIllegalFilenameRune(r rune) bool {
	switch r {
	case '<', '>', ':', '"', '/', '\\', '|', '?', '*':
		return true
	default:
		return false
	}
}

func EnsureUniqueFilename(base string, used map[string]int) string {
	if base == "" {
		base = "section"
	}

	n, exists := used[base]
	if !exists {
		used[base] = 1
		return base
	}

	n++
	used[base] = n
	return fmt.Sprintf("%s-%d", base, n)
}
