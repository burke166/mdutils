package output

import (
	"encoding/csv"
	"strconv"
	"strings"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

func RenderCsv(headings []markdown.Heading) (string, error) {
	var b strings.Builder
	w := csv.NewWriter(&b)

	if err := w.Write([]string{"level", "text"}); err != nil {
		return "", err
	}

	for _, h := range headings {
		if err := w.Write([]string{
			strconv.Itoa(h.Level),
			h.Text,
		}); err != nil {
			return "", err
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return "", err
	}

	return b.String(), nil
}
