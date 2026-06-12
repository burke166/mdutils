package mdmerge

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/computercodeblue/mdutils/internal/commands/mdsplit"
	"github.com/stretchr/testify/require"
)

func writeMarkdownFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}

func TestRunMergesExplicitFilesInOrder(t *testing.T) {
	dir := t.TempDir()
	first := writeMarkdownFile(t, dir, "01-intro.md", "# Intro\n\nIntro body.\n")
	second := writeMarkdownFile(t, dir, "02-rules.md", "# Rules\n\nRules body.\n")
	third := writeMarkdownFile(t, dir, "03-combat.md", "# Combat\n\nCombat body.\n")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{first, second, third}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	expected := MergeContents([]string{
		"# Intro\n\nIntro body.\n",
		"# Rules\n\nRules body.\n",
		"# Combat\n\nCombat body.\n",
	})
	require.Equal(t, expected, stdout.String())
}

func TestRunMergesDirectoryInFilenameOrder(t *testing.T) {
	dir := t.TempDir()
	writeMarkdownFile(t, dir, "introduction.md", "# Introduction\n\nWelcome.\n")
	writeMarkdownFile(t, dir, "character-creation.md", "# Character Creation\n\nMake a character.\n")
	writeMarkdownFile(t, dir, "combat.md", "# Combat\n\nCombat rules.\n")

	outPath := filepath.Join(t.TempDir(), "book-merged.md")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--out", outPath, dir}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	data, err := os.ReadFile(outPath)
	require.NoError(t, err)

	expected := MergeContents([]string{
		"# Character Creation\n\nMake a character.\n",
		"# Combat\n\nCombat rules.\n",
		"# Introduction\n\nWelcome.\n",
	})
	require.Equal(t, expected, string(data))
}

func TestRunWritesToStdoutWhenOutOmitted(t *testing.T) {
	dir := t.TempDir()
	first := writeMarkdownFile(t, dir, "a.md", "# A\n\nAlpha.\n")
	second := writeMarkdownFile(t, dir, "b.md", "# B\n\nBeta.\n")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{first, second}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)
	require.NotEmpty(t, stdout.String())
}

func TestRunWritesToOutputFile(t *testing.T) {
	dir := t.TempDir()
	path := writeMarkdownFile(t, dir, "only.md", "# Only\n\nBody.\n")
	outPath := filepath.Join(t.TempDir(), "merged.md")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--out", outPath, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	data, err := os.ReadFile(outPath)
	require.NoError(t, err)
	require.Equal(t, "# Only\n\nBody.\n", string(data))
}

func TestRunEnsuresSingleBlankLineBetweenFiles(t *testing.T) {
	dir := t.TempDir()
	first := writeMarkdownFile(t, dir, "a.md", "# A\n\nAlpha.\n\n\n")
	second := writeMarkdownFile(t, dir, "b.md", "\n# B\n\nBeta.\n")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{first, second}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)
	require.Equal(t, "# A\n\nAlpha.\n\n# B\n\nBeta.\n", stdout.String())
}

func TestRunDoesNotRecursivelyMergeSubdirectories(t *testing.T) {
	dir := t.TempDir()
	writeMarkdownFile(t, dir, "top.md", "# Top\n\nTop body.\n")
	subDir := filepath.Join(dir, "sub")
	require.NoError(t, os.Mkdir(subDir, 0755))
	writeMarkdownFile(t, subDir, "nested.md", "# Nested\n\nNested body.\n")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{dir}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)
	require.Equal(t, "# Top\n\nTop body.\n", stdout.String())
}

func TestRunMissingFile(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"missing.md"}, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 2, code)
}

func TestRunEmptyDirectory(t *testing.T) {
	dir := t.TempDir()

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{dir}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 1, code)
}

func TestRunMissingInputs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code, err := Run(nil, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 2, code)
	require.Contains(t, err.Error(), "missing input files or directory")
}

func TestRunSplitAndMergeRoundTrip(t *testing.T) {
	book := "# Introduction\n\nWelcome.\n\n# Character Creation\n\nMake a character.\n\n## Abilities\n\nAbility rules.\n\n# Combat\n\nCombat rules.\n"
	bookPath := writeMarkdownFile(t, t.TempDir(), "book.md", book)
	chaptersDir := filepath.Join(t.TempDir(), "chapters")

	var stdout, stderr bytes.Buffer
	code, err := mdsplit.Run([]string{"--out", chaptersDir, bookPath}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	outPath := filepath.Join(t.TempDir(), "book-merged.md")
	code, err = Run([]string{"--out", outPath, chaptersDir}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	data, err := os.ReadFile(outPath)
	require.NoError(t, err)

	expected := MergeContents([]string{
		"# Character Creation\n\nMake a character.\n\n## Abilities\n\nAbility rules.\n",
		"# Combat\n\nCombat rules.\n",
		"# Introduction\n\nWelcome.\n",
	})
	require.Equal(t, expected, string(data))
}
