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
	Line  int    `json:"line"`
}

func ExtractHeadings(source []byte) ([]Heading, error) {
	parser := goldmark.DefaultParser()
	doc := parser.Parse(text.NewReader(source))

	var headings []Heading
	searchFrom := 0

	err := ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		h, ok := n.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}

		line := 1
		nextSearchFrom := searchFrom

		if h.Lines().Len() > 0 {
			start := h.Lines().At(0).Start
			line = lineNumber(source, start)
			nextSearchFrom = h.Lines().At(0).Stop
		} else if offset, ok := findATXHeadingOffset(source, searchFrom, h.Level); ok {
			line = lineNumber(source, offset)
			nextSearchFrom = offset + 1
		}

		searchFrom = nextSearchFrom

		headings = append(headings, Heading{
			Level: h.Level,
			Text:  strings.TrimSpace(string(h.Text(source))),
			Line:  line,
		})

		return ast.WalkContinue, nil
	})

	return headings, err
}

func lineNumber(source []byte, offset int) int {
	if offset > len(source) {
		offset = len(source)
	}

	line := 1
	for i := 0; i < offset; i++ {
		if source[i] == '\n' {
			line++
		}
	}

	return line
}

func findATXHeadingOffset(source []byte, from int, level int) (int, bool) {
	if from < 0 || from >= len(source) {
		return 0, false
	}

	for i := from; i < len(source); i++ {
		if i > 0 && source[i-1] != '\n' {
			continue
		}
		if source[i] != '#' {
			continue
		}

		hashes := 0
		for j := i; j < len(source) && source[j] == '#'; j++ {
			hashes++
		}
		if hashes != level {
			continue
		}

		switch {
		case i+hashes >= len(source):
			return i, true
		case source[i+hashes] == ' ', source[i+hashes] == '\t', source[i+hashes] == '\n', source[i+hashes] == '\r':
			return i, true
		default:
			continue
		}
	}

	return 0, false
}
