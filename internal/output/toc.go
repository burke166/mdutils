package output

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

type TocOptions struct {
	MinLevel int
	MaxLevel int
	Ordered  bool
	NoLinks  bool
	NoIndent bool
}

func RenderToc(headings []markdown.Heading, opts TocOptions) string {
	slugs := AssignSlugs(headings)

	var b strings.Builder
	numbers := make([]int, 6)

	for i, h := range headings {
		level := h.Level
		if level < 1 {
			level = 1
		}
		if level > 6 {
			level = 6
		}

		var marker string
		if opts.Ordered {
			numbers[level-1]++

			for j := level; j < len(numbers); j++ {
				numbers[j] = 0
			}

			parts := make([]string, 0, level)
			for j := 0; j < level; j++ {
				if numbers[j] > 0 {
					parts = append(parts, strconv.Itoa(numbers[j]))
				}
			}

			marker = strings.Join(parts, ".") + "."
		}

		if h.Level < opts.MinLevel || h.Level > opts.MaxLevel {
			continue
		}

		if !opts.Ordered {
			marker = "-"
		}

		indent := ""
		if !opts.NoIndent {
			indent = strings.Repeat("  ", h.Level-1)
		}

		line := h.Text
		if !opts.NoLinks {
			line = fmt.Sprintf("[%s](#%s)", h.Text, slugs[i])
		}

		fmt.Fprintf(&b, "%s%s %s\n", indent, marker, line)
	}

	return b.String()
}
