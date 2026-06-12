package output

import (
	"fmt"
	"strings"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

func RenderMarkdownHeadings(headings []markdown.Heading) string {
	var b strings.Builder

	for _, h := range headings {
		level := h.Level
		if level < 1 {
			level = 1
		}
		if level > 6 {
			level = 6
		}

		fmt.Fprintf(&b, "%s %s\n", strings.Repeat("#", level), h.Text)
	}

	return b.String()
}
