package fileutil

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteFileAtomicallyPreservesPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permission bits are not preserved on Windows")
	}
	path := filepath.Join(t.TempDir(), "doc.md")
	require.NoError(t, os.WriteFile(path, []byte("original\n"), 0600))

	require.NoError(t, WriteFileAtomically(path, []byte("updated\n")))

	info, err := os.Stat(path)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0600), info.Mode().Perm())

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "updated\n", string(data))
}

func TestWriteFileAtomicallyCreatesNewFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new.md")

	require.NoError(t, WriteFileAtomically(path, []byte("content\n")))

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "content\n", string(data))
}
