package mdformat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func defaultFormatOptions() FormatOptions {
	return FormatOptions{
		TrimTrailingSpace:   true,
		EnsureFinalNewline:  true,
		NormalizeHeadings:   true,
		NormalizeLists:      true,
		PreserveFrontmatter: true,
	}
}

func TestFormatMarkdown(t *testing.T) {
	tests := []struct {
		name  string
		input string
		opts  FormatOptions
		want  string
	}{
		{
			name:  "normalizes heading spacing",
			input: "#Heading\n##  Heading Two\n###   Title ###\n",
			opts:  defaultFormatOptions(),
			want:  "# Heading\n## Heading Two\n### Title\n",
		},
		{
			name:  "removes trailing whitespace",
			input: "Line with spaces   \nAnother line\t \n",
			opts:  defaultFormatOptions(),
			want:  "Line with spaces\nAnother line\n",
		},
		{
			name:  "preserves hard line breaks when trim disabled",
			input: "First line  \nSecond line\n",
			opts: func() FormatOptions {
				opts := defaultFormatOptions()
				opts.TrimTrailingSpace = false
				return opts
			}(),
			want: "First line  \nSecond line\n",
		},
		{
			name:  "removes hard line break spaces when trim enabled",
			input: "First line  \nSecond line\n",
			opts:  defaultFormatOptions(),
			want:  "First line\nSecond line\n",
		},
		{
			name:  "ensures final newline",
			input: "# Title\n\nBody",
			opts:  defaultFormatOptions(),
			want:  "# Title\n\nBody\n",
		},
		{
			name:  "removes extra final newlines",
			input: "# Title\n\n\n",
			opts:  defaultFormatOptions(),
			want:  "# Title\n\n",
		},
		{
			name:  "preserves fenced code blocks",
			input: "# Real\n\n```\n# Fake Heading   \n* not a list\n```\n\n##Heading\n",
			opts:  defaultFormatOptions(),
			want:  "# Real\n\n```\n# Fake Heading   \n* not a list\n```\n\n## Heading\n",
		},
		{
			name:  "preserves YAML frontmatter",
			input: "---\ntitle: My Doc  \nextra: true\n---\n\n##Heading\n",
			opts:  defaultFormatOptions(),
			want:  "---\ntitle: My Doc  \nextra: true\n---\n\n## Heading\n",
		},
		{
			name:  "preserves frontmatter blank lines",
			input: "---\ntitle: Doc\n\nauthor: Me\n---\n\nBody\n",
			opts:  defaultFormatOptions(),
			want:  "---\ntitle: Doc\n\nauthor: Me\n---\n\nBody\n",
		},
		{
			name:  "normalizes unordered list markers",
			input: "* first\n+ second\n- third\n**not a list**\n",
			opts:  defaultFormatOptions(),
			want:  "- first\n- second\n- third\n**not a list**\n",
		},
		{
			name:  "preserves ordered list numbering",
			input: "1. first\n3. second\n10. third\n",
			opts:  defaultFormatOptions(),
			want:  "1. first\n3. second\n10. third\n",
		},
		{
			name:  "normalizes blank lines",
			input: "Paragraph one.\n\n\n\nParagraph two.\n",
			opts:  defaultFormatOptions(),
			want:  "Paragraph one.\n\nParagraph two.\n",
		},
		{
			name:  "preserves blank lines inside code fence",
			input: "```\nline one\n\n\nline two\n```\n",
			opts:  defaultFormatOptions(),
			want:  "```\nline one\n\n\nline two\n```\n",
		},
		{
			name:  "preserves indented code blocks",
			input: "    code line   \n    second line\n\nNormal text\n",
			opts:  defaultFormatOptions(),
			want:  "    code line   \n    second line\n\nNormal text\n",
		},
		{
			name:  "does not change inline code",
			input: "Use `* marker` and `## hash` in prose.\n",
			opts:  defaultFormatOptions(),
			want:  "Use `* marker` and `## hash` in prose.\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatMarkdown(tt.input, tt.opts)
			require.Equal(t, tt.want, got)
		})
	}
}
