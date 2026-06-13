package markdown

import "strings"

func SplitLines(content string) []string {
	if content == "" {
		return nil
	}
	return strings.SplitAfter(content, "\n")
}

func IsFenceLine(line string) (bool, byte) {
	trimmed := strings.TrimLeft(line, " \t")
	switch {
	case strings.HasPrefix(trimmed, "```"):
		return true, '`'
	case strings.HasPrefix(trimmed, "~~~"):
		return true, '~'
	default:
		return false, 0
	}
}

func IsBlankLine(line string) bool {
	return strings.TrimSpace(strings.TrimRight(line, "\r\n")) == ""
}

func TrailingNewline(line string) string {
	if strings.HasSuffix(line, "\r\n") {
		return "\r\n"
	}
	if strings.HasSuffix(line, "\n") {
		return "\n"
	}
	return ""
}
