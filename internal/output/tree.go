package output

import (
	"fmt"
	"strings"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

func RenderTree(headings []markdown.Heading) string {
	var b strings.Builder

	for i, h := range headings {
		prefix := treePrefix(headings, i)
		fmt.Fprintf(&b, "%s%s\n", prefix, h.Text)
	}

	return b.String()
}

func treePrefix(headings []markdown.Heading, index int) string {
	level := headings[index].Level

	if level <= 1 {
		return ""
	}

	var b strings.Builder

	for ancestorLevel := 2; ancestorLevel < level; ancestorLevel++ {
		if hasLaterSiblingAtLevel(headings, index, ancestorLevel) {
			b.WriteString("│   ")
		} else {
			b.WriteString("    ")
		}
	}

	if hasLaterSiblingAtLevel(headings, index, level) {
		b.WriteString("├── ")
	} else {
		b.WriteString("└── ")
	}

	return b.String()
}

func hasLaterSiblingAtLevel(headings []markdown.Heading, index int, level int) bool {
	for i := index + 1; i < len(headings); i++ {
		if headings[i].Level < level {
			return false
		}

		if headings[i].Level == level {
			return true
		}
	}

	return false
}
