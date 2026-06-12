package output

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

func RenderNumbered(headings []markdown.Heading) string {
	var b strings.Builder
	numbers := make([]int, 6)

	for _, h := range headings {
		level := h.Level

		if level < 1 || level > 6 {
			continue
		}

		numbers[level-1]++

		for i := level; i < len(numbers); i++ {
			numbers[i] = 0
		}

		parts := make([]string, 0, level)

		for i := 0; i < level; i++ {
			if numbers[i] > 0 {
				parts = append(parts, strconv.Itoa(numbers[i]))
			}
		}

		indent := strings.Repeat("  ", level-1)
		fmt.Fprintf(&b, "%s%s. %s\n", indent, strings.Join(parts, "."), h.Text)
	}

	return b.String()
}
