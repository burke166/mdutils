package mdformat

import (
	"strings"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

type FormatOptions struct {
	TrimTrailingSpace   bool
	EnsureFinalNewline  bool
	NormalizeHeadings   bool
	NormalizeLists      bool
	PreserveFrontmatter bool
}

func FormatMarkdown(content string, opts FormatOptions) string {
	lines := markdown.SplitLines(content)
	if len(lines) == 0 {
		return ensureFinalNewline("", opts.EnsureFinalNewline)
	}

	var out []string
	bodyStart := 0

	if opts.PreserveFrontmatter {
		if count, ok := markdown.FrontmatterBounds(lines); ok {
			out = append(out, lines[:count]...)
			bodyStart = count
		}
	}

	var prevBlank bool

	inFence := false
	var fenceChar byte
	inIndentedCode := false

	for _, line := range lines[bodyStart:] {
		isFence, char := markdown.IsFenceLine(line)
		if isFence {
			if !inFence {
				inFence = true
				fenceChar = char
			} else if char == fenceChar {
				inFence = false
				fenceChar = 0
			}
			out = append(out, line)
			inIndentedCode = false
			prevBlank = false
			continue
		}

		if inFence {
			out = append(out, line)
			prevBlank = false
			continue
		}

		if inIndentedCode {
			if markdown.IsBlankLine(line) || isIndentedCodeLine(line) {
				out = append(out, line)
				prevBlank = markdown.IsBlankLine(line)
				continue
			}
			inIndentedCode = false
		}

		if isIndentedCodeLine(line) {
			inIndentedCode = true
			out = append(out, line)
			prevBlank = false
			continue
		}

		formatted := formatNormalLine(line, opts)

		if markdown.IsBlankLine(formatted) {
			if prevBlank {
				continue
			}
			prevBlank = true
		} else {
			prevBlank = false
		}

		out = append(out, formatted)
	}

	result := strings.Join(out, "")
	return ensureFinalNewline(result, opts.EnsureFinalNewline)
}

func formatNormalLine(line string, opts FormatOptions) string {
	formatted := line
	if opts.NormalizeHeadings {
		formatted = normalizeATXHeading(formatted)
	}
	if opts.NormalizeLists {
		formatted = normalizeUnorderedListMarker(formatted)
	}
	if opts.TrimTrailingSpace {
		formatted = trimTrailingWhitespace(formatted)
	}
	return formatted
}

func normalizeATXHeading(line string) string {
	leading := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
	trimmed := strings.TrimLeft(line, " \t")
	if len(trimmed) == 0 || trimmed[0] != '#' {
		return line
	}

	level := 0
	for level < len(trimmed) && trimmed[level] == '#' {
		level++
	}
	if level == 0 || level > 6 {
		return line
	}

	rest := trimmed[level:]
	newline := markdown.TrailingNewline(line)

	body := strings.TrimRight(rest, "\r\n")
	body = strings.TrimSpace(body)
	body = strings.TrimRight(body, " #")
	body = strings.TrimSpace(body)

	hashes := strings.Repeat("#", level)
	if body == "" {
		return leading + hashes + newline
	}

	return leading + hashes + " " + body + newline
}

func normalizeUnorderedListMarker(line string) string {
	leading := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
	trimmed := strings.TrimLeft(line, " \t")
	if len(trimmed) == 0 {
		return line
	}

	switch trimmed[0] {
	case '*', '+':
	default:
		return line
	}

	if len(trimmed) == 1 {
		return line
	}

	if trimmed[0] == '*' && len(trimmed) > 1 && trimmed[1] == '*' {
		return line
	}

	if trimmed[1] != ' ' && trimmed[1] != '\t' {
		return line
	}

	rest := trimmed[1:]
	if len(rest) > 0 && (rest[0] == ' ' || rest[0] == '\t') {
		rest = rest[1:]
	}

	newline := markdown.TrailingNewline(line)
	return leading + "- " + rest + newline
}

func trimTrailingWhitespace(line string) string {
	newline := markdown.TrailingNewline(line)
	body := strings.TrimRight(line, "\r\n")
	body = strings.TrimRight(body, " \t")
	return body + newline
}

func ensureFinalNewline(content string, enabled bool) string {
	if !enabled {
		return content
	}
	if content == "" {
		return "\n"
	}
	if strings.HasSuffix(content, "\n") {
		return content
	}
	return content + "\n"
}

func isIndentedCodeLine(line string) bool {
	if strings.HasPrefix(line, "\t") {
		return true
	}

	leading := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
	return len(leading) >= 4
}
