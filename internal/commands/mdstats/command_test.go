package mdstats

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunSingleFile(t *testing.T) {
	path := testdataPath(t, "simple.md")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	output := stdout.String()
	require.Contains(t, output, "simple.md")
	require.Contains(t, output, "Lines:")
	require.Contains(t, output, "Headings:")
	require.Contains(t, output, "Markdown:")
}

func TestRunFolderFlat(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "b.md"), []byte("# B\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.md"), []byte("# A\n"), 0644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "nested"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "nested", "c.md"), []byte("# C\n"), 0644))

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--no-recursive", dir}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	output := stdout.String()
	require.Contains(t, output, filepath.Join(dir, "a.md"))
	require.Contains(t, output, filepath.Join(dir, "b.md"))
	require.NotContains(t, output, filepath.Join(dir, "nested", "c.md"))
}

func TestRunFolderRecursive(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.md"), []byte("# A\n"), 0644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "nested"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "nested", "b.md"), []byte("# B\n"), 0644))

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{dir}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	output := stdout.String()
	require.Contains(t, output, filepath.Join(dir, "a.md"))
	require.Contains(t, output, filepath.Join(dir, "nested", "b.md"))
}

func TestRunIgnoresVendorFolders(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.md"), []byte("# A\n"), 0644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "vendor"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "vendor", "b.md"), []byte("# B\n"), 0644))

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{dir}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)
	require.NotContains(t, stdout.String(), filepath.Join(dir, "vendor", "b.md"))
}

func TestRunExcludeGlobs(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Readme\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "CHANGELOG.md"), []byte("# Change\n"), 0644))
	require.NoError(t, os.MkdirAll(filepath.Join(dir, "docs", "generated"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "docs", "generated", "out.md"), []byte("# Out\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "docs", "manual.md"), []byte("# Manual\n"), 0644))

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{
		"--exclude", "CHANGELOG.md",
		"--exclude", "docs/generated/**",
		dir,
	}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	output := stdout.String()
	require.Contains(t, output, "README.md")
	require.Contains(t, output, filepath.Join("docs", "manual.md"))
	require.NotContains(t, output, "CHANGELOG.md")
	require.NotContains(t, output, filepath.Join("docs", "generated", "out.md"))
}

func TestRunJSONOutput(t *testing.T) {
	path := testdataPath(t, "simple.md")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--json", path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	var result AnalysisResult
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &result))
	require.Len(t, result.Files, 1)
	require.Equal(t, 1, result.Summary.FileCount)
	require.Equal(t, result.Files[0].Words, result.Summary.Words)
	require.Contains(t, stdout.String(), `"fileSizeBytes"`)
	require.Contains(t, stdout.String(), `"readingTimeMinutes"`)
	require.Contains(t, stdout.String(), `"bulletItems"`)
}

func TestRunCSVOutput(t *testing.T) {
	path := testdataPath(t, "simple.md")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--csv", path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	output := strings.TrimSpace(stdout.String())
	lines := strings.Split(output, "\n")
	require.GreaterOrEqual(t, len(lines), 2)
	require.Equal(t, "path,fileSizeBytes,lines,blankLines,characters,words,paragraphs,sentences,readingTimeMinutes,h1,h2,h3,h4,h5,h6,totalHeadings,maxHeadingDepth,lists,bulletItems,numberedItems,taskItems,blockQuoteLines,codeBlocks,inlineCodeSpans,tables,links,images,footnotes,horizontalRules,frontmatterDetected,frontmatterLines", lines[0])
	require.True(t, strings.Contains(lines[1], "simple.md"))
}

func TestRunSummaryOutput(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.md"), []byte("one two three\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "b.md"), []byte("four five six\n"), 0644))

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--summary", "--per-file=false", dir}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	output := stdout.String()
	require.Contains(t, output, "Summary")
	require.Contains(t, output, "Files:")
	require.NotContains(t, output, "a.md")
}

func TestRunStableFileSorting(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "z.md"), []byte("# Z\n"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.md"), []byte("# A\n"), 0644))

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--json", dir}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	var result AnalysisResult
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &result))
	require.Len(t, result.Files, 2)
	require.True(t, strings.HasSuffix(result.Files[0].Path, "a.md"))
	require.True(t, strings.HasSuffix(result.Files[1].Path, "z.md"))
}

func TestRunUnreadableFileExitCode(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.md")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{path}, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 3, code)
}

func TestRunPartialReadErrors(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "good.md")
	require.NoError(t, os.WriteFile(good, []byte("# Good\n"), 0644))

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{dir, good}, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 3, code)
}

func TestRunMultipleFormatsRejected(t *testing.T) {
	path := testdataPath(t, "simple.md")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--json", "--csv", path}, &stdout, &stderr)
	require.EqualError(t, err, "choose only one output format")
	require.Equal(t, 3, code)
}

func TestRunMissingInputArgument(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code, err := Run(nil, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 3, code)
	require.Contains(t, err.Error(), "missing input file or directory")
}
