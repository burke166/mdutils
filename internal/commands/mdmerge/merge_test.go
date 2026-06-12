package mdmerge

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeContents(t *testing.T) {
	merged := MergeContents([]string{
		"# One\n\nFirst.\n",
		"# Two\n\nSecond.\n",
	})
	require.Equal(t, "# One\n\nFirst.\n\n# Two\n\nSecond.\n", merged)
}

func TestMergeContentsSingleBlankLineBetweenFiles(t *testing.T) {
	merged := MergeContents([]string{
		"# One\n\nFirst.\n\n\n",
		"\n# Two\n\nSecond.\n",
	})
	require.Equal(t, "# One\n\nFirst.\n\n# Two\n\nSecond.\n", merged)
}

func TestMergeContentsEmpty(t *testing.T) {
	require.Equal(t, "", MergeContents(nil))
}
