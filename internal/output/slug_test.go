package output

import (
	"testing"

	"github.com/computercodeblue/mdutils/internal/markdown"
	"github.com/stretchr/testify/require"
)

func TestSlug(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "simple",
			text:     "Usage",
			expected: "usage",
		},
		{
			name:     "trim spaces",
			text:     "  Markdown output  ",
			expected: "markdown-output",
		},
		{
			name:     "punctuation",
			text:     "Usage!",
			expected: "usage",
		},
		{
			name:     "mixed punctuation and spaces",
			text:     "Hello, World!",
			expected: "hello-world",
		},
		{
			name:     "preserve numbers",
			text:     "Section 2",
			expected: "section-2",
		},
		{
			name:     "collapse hyphens",
			text:     "foo   bar",
			expected: "foo-bar",
		},
		{
			name:     "only punctuation",
			text:     "!!!",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, Slug(tt.text))
		})
	}
}

func TestAssignSlugs(t *testing.T) {
	headings := []markdown.Heading{
		{Level: 1, Text: "Usage"},
		{Level: 2, Text: "Usage"},
		{Level: 2, Text: "Usage!"},
	}

	require.Equal(t, []string{"usage", "usage-1", "usage-2"}, AssignSlugs(headings))
}

func TestAssignSlugsEmptyBase(t *testing.T) {
	headings := []markdown.Heading{
		{Level: 1, Text: "!!!"},
		{Level: 2, Text: "!!!"},
	}

	require.Equal(t, []string{"", "-1"}, AssignSlugs(headings))
}
