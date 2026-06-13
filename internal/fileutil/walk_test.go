package fileutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatchGlobSimple(t *testing.T) {
	require.True(t, matchGlob("CHANGELOG.md", "CHANGELOG.md"))
	require.True(t, matchGlob("*.md", "notes.md"))
	require.False(t, matchGlob("*.md", "notes.txt"))
}

func TestMatchGlobDoubleStar(t *testing.T) {
	require.True(t, matchGlob("docs/generated/**", "docs/generated/output.md"))
	require.True(t, matchGlob("docs/generated/**", "docs/generated/nested/output.md"))
	require.False(t, matchGlob("docs/generated/**", "output.md"))
	require.False(t, matchGlob("docs/generated/**", "docs/manual.md"))
}

func TestFilterExcluded(t *testing.T) {
	files := []string{
		"README.md",
		"CHANGELOG.md",
		"docs/generated/output.md",
		"docs/manual.md",
	}

	filtered := FilterExcluded(files, ".", []string{
		"CHANGELOG.md",
		"docs/generated/**",
	})

	require.Equal(t, []string{"README.md", "docs/manual.md"}, filtered)
}

func TestCollectMarkdownFilesSingleFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doc.md")
	require.NoError(t, os.WriteFile(path, []byte("# Title\n"), 0644))

	files, err := CollectMarkdownFiles(path, true)
	require.NoError(t, err)
	require.Equal(t, []string{path}, files)
}

func TestCollectMarkdownFilesFlatDirectory(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.md"), []byte("# A\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "b.markdown"), []byte("# B\n"), 0644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "nested"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "nested", "c.md"), []byte("# C\n"), 0644))

	files, err := CollectMarkdownFiles(dir, false)
	require.NoError(t, err)
	require.Len(t, files, 2)
}

func TestCollectMarkdownFilesRecursiveDirectory(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.md"), []byte("# A\n"), 0644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "nested"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "nested", "b.md"), []byte("# B\n"), 0644))

	files, err := CollectMarkdownFiles(dir, true)
	require.NoError(t, err)
	require.Len(t, files, 2)
}

func TestCollectMarkdownFilesIgnoresVendorFolders(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.md"), []byte("# A\n"), 0644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "node_modules"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "node_modules", "b.md"), []byte("# B\n"), 0644))

	files, err := CollectMarkdownFiles(dir, true)
	require.NoError(t, err)
	require.Len(t, files, 1)
}

func TestCollectMarkdownFilesMissingDirectory(t *testing.T) {
	_, err := CollectMarkdownFiles(filepath.Join(t.TempDir(), "missing"), true)
	require.Error(t, err)
}

func TestCollectMarkdownFilesSorted(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "z.md"), []byte("# Z\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.md"), []byte("# A\n"), 0644))

	files, err := CollectMarkdownFiles(dir, false)
	require.NoError(t, err)
	require.Equal(t, []string{filepath.Join(dir, "a.md"), filepath.Join(dir, "z.md")}, files)
}

func TestShouldIgnoreDir(t *testing.T) {
	require.True(t, shouldIgnoreDir("node_modules"))
	require.False(t, shouldIgnoreDir("docs"))
}
