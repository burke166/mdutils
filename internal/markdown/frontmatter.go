package markdown

import "strings"

func IsFrontmatterDelimiter(line string) bool {
	return strings.TrimSpace(strings.TrimRight(line, "\r\n")) == "---"
}

// FrontmatterBounds returns the number of lines occupied by YAML frontmatter
// at the start of the document, including opening and closing delimiters.
// When frontmatter is not present, ok is false and lineCount is zero.
func FrontmatterBounds(lines []string) (lineCount int, ok bool) {
	if len(lines) == 0 || !IsFrontmatterDelimiter(lines[0]) {
		return 0, false
	}

	for i := 1; i < len(lines); i++ {
		if IsFrontmatterDelimiter(lines[i]) {
			return i + 1, true
		}
	}

	return 0, false
}
