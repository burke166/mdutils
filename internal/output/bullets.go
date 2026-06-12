package output

import (
	"fmt"
	"strings"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

func RenderBullets(headings []markdown.Heading) string {
	var b strings.Builder

	for _, h := range headings {
		indent := strings.Repeat("  ", h.Level-1)
		fmt.Fprintf(&b, "%s- %s\n", indent, h.Text)
	}

	return b.String()
}
