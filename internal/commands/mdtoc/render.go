package mdtoc

import (
	"github.com/computercodeblue/mdutils/internal/markdown"
	"github.com/computercodeblue/mdutils/internal/output"
)

type RenderOptions = output.TocOptions

func Render(headings []markdown.Heading, opts RenderOptions) string {
	return output.RenderToc(headings, opts)
}
