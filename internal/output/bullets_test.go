package output

import (
	"strings"
	"testing"

	"github.com/computercodeblue/mdutils/internal/markdown"
	"github.com/stretchr/testify/require"
)

func TestBullets(t *testing.T) {
	tests := []struct {
		name     string
		headings []markdown.Heading
		expected string
	}{
		{
			name:     "empty",
			headings: emptyHeadings,
			expected: "",
		},
		{
			name:     "single",
			headings: singleHeading,
			expected: "- README",
		},
		{
			name:     "nested",
			headings: testHeadings,
			expected: strings.TrimSpace(`
- Adventure Wargame
  - Character Creation
    - Attributes
  - Equipment
			`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := strings.TrimSpace(RenderBullets(tt.headings))
			require.Equal(t, tt.expected, actual)
		})
	}
}
