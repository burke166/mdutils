package output

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

func Slug(text string) string {
	s := strings.ToLower(strings.TrimSpace(text))

	var b strings.Builder
	lastWasSpace := false

	for _, r := range s {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastWasSpace = false
		case unicode.IsSpace(r):
			if !lastWasSpace && b.Len() > 0 {
				b.WriteByte(' ')
				lastWasSpace = true
			}
		}
	}

	slug := strings.ReplaceAll(strings.TrimSpace(b.String()), " ", "-")

	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	return strings.Trim(slug, "-")
}

func AssignSlugs(headings []markdown.Heading) []string {
	slugs := make([]string, len(headings))
	counts := make(map[string]int)

	for i, h := range headings {
		base := Slug(h.Text)
		n := counts[base]
		counts[base] = n + 1

		if n == 0 {
			slugs[i] = base
		} else {
			slugs[i] = fmt.Sprintf("%s-%d", base, n)
		}
	}

	return slugs
}
