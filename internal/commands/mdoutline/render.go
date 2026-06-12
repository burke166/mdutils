package mdoutline

import (
	"fmt"

	"github.com/computercodeblue/mdutils/internal/markdown"
	"github.com/computercodeblue/mdutils/internal/output"
)

func Render(headings []markdown.Heading, format string) (string, error) {
	switch format {
	case "bullets":
		return output.RenderBullets(headings), nil
	case "tree":
		return output.RenderTree(headings), nil
	case "numbered":
		return output.RenderNumbered(headings), nil
	case "json":
		return output.RenderJson(headings)
	case "csv":
		return output.RenderCsv(headings)
	case "headings":
		return output.RenderMarkdownHeadings(headings), nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}
