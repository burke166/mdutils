package output

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/computercodeblue/mdutils/internal/markdown"
	"github.com/stretchr/testify/require"
)

func TestRenderToc(t *testing.T) {
	tests := []struct {
		name     string
		headings []markdown.Heading
		opts     TocOptions
		expected string
	}{
		{
			name:     "unordered links",
			headings: tocHeadings,
			opts:     TocOptions{MinLevel: 1, MaxLevel: 6},
			expected: strings.TrimSpace(`
- [Install](#install)
- [Usage](#usage)
  - [Markdown output](#markdown-output)
  - [JSON output](#json-output)
- [License](#license)
			`),
		},
		{
			name:     "ordered links",
			headings: tocHeadings,
			opts:     TocOptions{MinLevel: 1, MaxLevel: 6, Ordered: true},
			expected: strings.TrimSpace(`
1. [Install](#install)
2. [Usage](#usage)
  2.1. [Markdown output](#markdown-output)
  2.2. [JSON output](#json-output)
3. [License](#license)
			`),
		},
		{
			name:     "unordered no links",
			headings: tocHeadings,
			opts:     TocOptions{MinLevel: 1, MaxLevel: 6, NoLinks: true},
			expected: strings.TrimSpace(`
- Install
- Usage
  - Markdown output
  - JSON output
- License
			`),
		},
		{
			name:     "ordered no links",
			headings: tocHeadings,
			opts:     TocOptions{MinLevel: 1, MaxLevel: 6, Ordered: true, NoLinks: true},
			expected: strings.TrimSpace(`
1. Install
2. Usage
  2.1. Markdown output
  2.2. JSON output
3. License
			`),
		},
		{
			name:     "unordered no indent",
			headings: tocHeadings,
			opts:     TocOptions{MinLevel: 1, MaxLevel: 6, NoIndent: true},
			expected: strings.TrimSpace(`
- [Install](#install)
- [Usage](#usage)
- [Markdown output](#markdown-output)
- [JSON output](#json-output)
- [License](#license)
			`),
		},
		{
			name:     "ordered no indent",
			headings: tocHeadings,
			opts:     TocOptions{MinLevel: 1, MaxLevel: 6, Ordered: true, NoIndent: true},
			expected: strings.TrimSpace(`
1. [Install](#install)
2. [Usage](#usage)
2.1. [Markdown output](#markdown-output)
2.2. [JSON output](#json-output)
3. [License](#license)
			`),
		},
		{
			name: "level filtering",
			headings: []markdown.Heading{
				{Level: 1, Text: "Title"},
				{Level: 2, Text: "Section"},
				{Level: 3, Text: "Subsection"},
				{Level: 4, Text: "Detail"},
			},
			opts: TocOptions{MinLevel: 2, MaxLevel: 3},
			expected: strings.TrimSpace(`
  - [Section](#section)
    - [Subsection](#subsection)
			`),
		},
		{
			name: "nested levels",
			headings: []markdown.Heading{
				{Level: 1, Text: "Adventure Wargame"},
				{Level: 2, Text: "Character Creation"},
				{Level: 3, Text: "Attributes"},
				{Level: 2, Text: "Equipment"},
			},
			opts: TocOptions{MinLevel: 1, MaxLevel: 6},
			expected: strings.TrimSpace(`
- [Adventure Wargame](#adventure-wargame)
  - [Character Creation](#character-creation)
    - [Attributes](#attributes)
  - [Equipment](#equipment)
			`),
		},
		{
			name: "duplicate slugs",
			headings: []markdown.Heading{
				{Level: 1, Text: "Usage"},
				{Level: 2, Text: "Usage"},
				{Level: 2, Text: "Usage!"},
			},
			opts: TocOptions{MinLevel: 1, MaxLevel: 6},
			expected: strings.TrimSpace(`
- [Usage](#usage)
  - [Usage](#usage-1)
  - [Usage!](#usage-2)
			`),
		},
		{
			name: "punctuation in headings",
			headings: []markdown.Heading{
				{Level: 1, Text: "Hello, World!"},
				{Level: 2, Text: "What's next?"},
			},
			opts: TocOptions{MinLevel: 1, MaxLevel: 6},
			expected: strings.TrimSpace(`
- [Hello, World!](#hello-world)
  - [What's next?](#whats-next)
			`),
		},
		{
			name: "hierarchical ordered",
			headings: []markdown.Heading{
				{Level: 1, Text: "First"},
				{Level: 2, Text: "First sub-level"},
				{Level: 2, Text: "Second sub-level"},
				{Level: 3, Text: "First sub-sub-level"},
				{Level: 2, Text: "Third sub-level"},
				{Level: 1, Text: "Second"},
			},
			opts: TocOptions{MinLevel: 1, MaxLevel: 6, Ordered: true, NoLinks: true},
			expected: strings.TrimSpace(`
1. First
  1.1. First sub-level
  1.2. Second sub-level
    1.2.1. First sub-sub-level
  1.3. Third sub-level
2. Second
			`),
		},
		{
			name:     "empty",
			headings: nil,
			opts:     TocOptions{MinLevel: 1, MaxLevel: 6},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := strings.TrimSpace(RenderToc(tt.headings, tt.opts))
			require.Equal(t, tt.expected, actual)
		})
	}
}

func TestRenderTocNestedDocument(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "nested.md")

	source, err := os.ReadFile(path)
	require.NoError(t, err)

	headings, err := markdown.ExtractHeadings(source)
	require.NoError(t, err)

	output := strings.TrimSpace(RenderToc(headings, TocOptions{
		MinLevel: 1,
		MaxLevel: 3,
	}))

	require.Contains(t, output, "- [Adventure Wargame](#adventure-wargame)")
	require.Contains(t, output, "  - [Character Creation](#character-creation)")
	require.Contains(t, output, "    - [Attributes](#attributes)")
	require.NotContains(t, output, "Body")
}
