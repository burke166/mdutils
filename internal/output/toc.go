package output

import (
	"fmt"
	"strings"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

type TocOptions struct {
	MinLevel int
	MaxLevel int
	Ordered  bool
	NoLinks  bool
}

func RenderToc(headings []markdown.Heading, opts TocOptions) string {
	slugs := AssignSlugs(headings)

	var b strings.Builder

	for i, h := range headings {
		if h.Level < opts.MinLevel || h.Level > opts.MaxLevel {
			continue
		}

		indent := strings.Repeat("  ", h.Level-1)

		marker := "-"
		if opts.Ordered {
			marker = "1."
		}

		line := h.Text
		if !opts.NoLinks {
			line = fmt.Sprintf("[%s](#%s)", h.Text, slugs[i])
		}

		fmt.Fprintf(&b, "%s%s %s\n", indent, marker, line)
	}

	return b.String()
}
