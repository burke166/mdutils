package markdown

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImportMarkdownHeadingsFromFiles(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     []Heading
	}{
		{
			name:     "simple document",
			fileName: "simple.md",
			want: []Heading{
				{Level: 1, Text: "Simple Document"},
				{Level: 2, Text: "Getting Started"},
				{Level: 3, Text: "Installation"},
				{Level: 3, Text: "Usage"},
				{Level: 2, Text: "Configuration"},
				{Level: 3, Text: "Output"},
				{Level: 2, Text: "Conclusion"},
			},
		},
		{
			name:     "ignores headings inside code blocks",
			fileName: "codeblocks.md",
			want: []Heading{
				{Level: 1, Text: "Code Block Test"},
				{Level: 2, Text: "C#"},
				{Level: 2, Text: "Bash"},
				{Level: 2, Text: "Markdown Example"},
				{Level: 2, Text: "Real Heading"},
			},
		},
		{
			name:     "ignores YAML front matter",
			fileName: "frontmatter.md",
			want: []Heading{
				{Level: 1, Text: "Front Matter Test"},
				{Level: 2, Text: "Section One"},
				{Level: 2, Text: "Section Two"},
				{Level: 3, Text: "Details"},
				{Level: 2, Text: "Final Notes"},
			},
		},
		{
			name:     "nested headings",
			fileName: "nested.md",
			want: []Heading{
				{Level: 1, Text: "Adventure Wargame"},
				{Level: 2, Text: "Character Creation"},
				{Level: 3, Text: "Attributes"},
				{Level: 4, Text: "Body"},
				{Level: 4, Text: "Agility"},
				{Level: 4, Text: "Intellect"},
				{Level: 4, Text: "Presence"},
				{Level: 3, Text: "Skills"},
				{Level: 4, Text: "Combat"},
				{Level: 5, Text: "Firearms"},
				{Level: 5, Text: "Heavy Weapons"},
				{Level: 4, Text: "Technical"},
				{Level: 5, Text: "Electronics"},
				{Level: 5, Text: "Mechanics"},
				{Level: 2, Text: "Equipment"},
				{Level: 3, Text: "Weapons"},
				{Level: 3, Text: "Armor"},
				{Level: 2, Text: "Combat"},
				{Level: 3, Text: "Initiative"},
				{Level: 3, Text: "Actions"},
				{Level: 3, Text: "Damage"},
				{Level: 2, Text: "Appendix"},
				{Level: 3, Text: "Tables"},
				{Level: 3, Text: "Index"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join("..", "..", "testdata", tt.fileName)

			source, err := os.ReadFile(path)
			require.NoError(t, err)

			got, err := ExtractHeadings(source)
			require.NoError(t, err)

			require.Equal(t, tt.want, got)
		})
	}
}
