package mdoutline

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func testdataPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join("..", "..", "..", "testdata", name)
}

func TestRunSimpleDocumentDefaultBullets(t *testing.T) {
	path := testdataPath(t, "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)

	expected := strings.TrimSpace(`
- Simple Document
  - Getting Started
    - Installation
    - Usage
  - Configuration
    - Output
  - Conclusion
	`)

	require.Equal(t, expected, strings.TrimSpace(stdout.String()))
}

func TestRunCodeBlocksDocument(t *testing.T) {
	path := testdataPath(t, "codeblocks.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)

	output := stdout.String()
	require.Contains(t, output, "- Code Block Test")
	require.Contains(t, output, "- Real Heading")
	require.NotContains(t, output, "Fake Heading")
}

func TestRunFrontMatterDocument(t *testing.T) {
	path := testdataPath(t, "frontmatter.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)

	expected := strings.TrimSpace(`
- Front Matter Test
  - Section One
  - Section Two
    - Details
  - Final Notes
	`)

	require.Equal(t, expected, strings.TrimSpace(stdout.String()))
}

func TestRunMissingFile(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"does-not-exist.md"}, &stdout, &stderr)
	require.Error(t, err)
}

func TestRunMissingInputArgument(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(nil, &stdout, &stderr)
	require.Error(t, err)
	require.Contains(t, err.Error(), "missing input Markdown file")
}

func TestRunTreeFormat(t *testing.T) {
	path := testdataPath(t, "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--tree", path}, &stdout, &stderr)
	require.NoError(t, err)

	output := stdout.String()
	require.Contains(t, output, "Simple Document")
	require.Contains(t, output, "Getting Started")
	require.Contains(t, output, "├──")
}

func TestRunNumberedFormat(t *testing.T) {
	path := testdataPath(t, "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--numbered", path}, &stdout, &stderr)
	require.NoError(t, err)

	output := stdout.String()
	require.Contains(t, output, "1. Simple Document")
	require.Contains(t, output, "1.1. Getting Started")
}

func TestRunJsonFormat(t *testing.T) {
	path := testdataPath(t, "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--json", path}, &stdout, &stderr)
	require.NoError(t, err)

	output := stdout.String()
	require.Contains(t, output, `"level": 1`)
	require.Contains(t, output, `"text": "Simple Document"`)
}

func TestRunCsvFormat(t *testing.T) {
	path := testdataPath(t, "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--csv", path}, &stdout, &stderr)
	require.NoError(t, err)

	output := strings.TrimSpace(stdout.String())
	require.True(t, strings.HasPrefix(output, "level,text"))
	require.Contains(t, output, "1,Simple Document")
	require.Contains(t, output, "2,Getting Started")
}

func TestRunHeadingsFormat(t *testing.T) {
	path := testdataPath(t, "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--headings", path}, &stdout, &stderr)
	require.NoError(t, err)

	output := stdout.String()
	require.Contains(t, output, "# Simple Document")
	require.Contains(t, output, "## Getting Started")
	require.Contains(t, output, "### Installation")
}

func TestRunBulletsFormat(t *testing.T) {
	path := testdataPath(t, "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--bullets", path}, &stdout, &stderr)
	require.NoError(t, err)

	expected := strings.TrimSpace(`
- Simple Document
  - Getting Started
    - Installation
    - Usage
  - Configuration
    - Output
  - Conclusion
	`)

	require.Equal(t, expected, strings.TrimSpace(stdout.String()))
}

func TestRunMultipleFormats(t *testing.T) {
	path := testdataPath(t, "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--bullets", "--tree", path}, &stdout, &stderr)
	require.EqualError(t, err, "choose only one output format")
}

func TestRunOutputFile(t *testing.T) {
	path := testdataPath(t, "simple.md")
	outPath := filepath.Join(t.TempDir(), "outline.txt")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"-o", outPath, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Empty(t, stdout.String())

	written, err := os.ReadFile(outPath)
	require.NoError(t, err)
	require.Contains(t, string(written), "- Simple Document")
}

func TestRunOutputFileLongFlag(t *testing.T) {
	path := testdataPath(t, "simple.md")
	outPath := filepath.Join(t.TempDir(), "outline.txt")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--output", outPath, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Empty(t, stdout.String())

	written, err := os.ReadFile(outPath)
	require.NoError(t, err)
	require.Contains(t, string(written), "- Simple Document")
}
