package mdmerge

import (
	"errors"
	"os"
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

func TestCollectExitCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "nil", err: nil, want: 0},
		{name: "not exist", err: os.ErrNotExist, want: 2},
		{name: "path error", err: &os.PathError{Op: "read", Path: "dir", Err: os.ErrPermission}, want: 2},
		{name: "empty directory", err: errors.New("no Markdown files found in directory: docs"), want: 1},
		{name: "not markdown", err: errors.New("not a Markdown file: notes.txt"), want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, CollectExitCode(tt.err))
		})
	}
}
