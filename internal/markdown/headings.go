package markdown

import (
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type Heading struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
}

func ExtractHeadings(source []byte) ([]Heading, error) {
	parser := goldmark.DefaultParser()
	doc := parser.Parse(text.NewReader(source))

	var headings []Heading

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		h, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}

		headings = append(headings, Heading{
			Level: h.Level,
			Text:  strings.TrimSpace(string(h.Text(source))),
		})

		return ast.WalkContinue, nil
	})

	return headings, err
}
