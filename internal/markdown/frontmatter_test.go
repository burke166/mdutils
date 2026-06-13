package markdown

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFrontmatterBounds(t *testing.T) {
	lines := SplitLines("---\ntitle: Test\n---\n\n# Body\n")
	count, ok := FrontmatterBounds(lines)
	require.True(t, ok)
	require.Equal(t, 3, count)
}

func TestFrontmatterBoundsMissingClosingDelimiter(t *testing.T) {
	lines := SplitLines("---\ntitle: Test\n# Body\n")
	_, ok := FrontmatterBounds(lines)
	require.False(t, ok)
}

func TestFrontmatterBoundsNoFrontmatter(t *testing.T) {
	lines := SplitLines("# Title\n")
	_, ok := FrontmatterBounds(lines)
	require.False(t, ok)
}
