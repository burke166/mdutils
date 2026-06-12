package mdcheck

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunValidDocument(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)
	require.Empty(t, stdout.String())
}

func TestRunMissingH1(t *testing.T) {
	path := writeTempMarkdown(t, "## Section\n\n### Subsection\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 1, code)
	require.Contains(t, stdout.String(), "[missing-h1]")
}

func TestRunMultipleH1(t *testing.T) {
	path := writeTempMarkdown(t, "# First\n\n# Second\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 1, code)
	require.Contains(t, stdout.String(), "[multiple-h1]")
	require.Contains(t, stdout.String(), "Second")
}

func TestRunAllowMultipleH1(t *testing.T) {
	path := writeTempMarkdown(t, "# First\n\n# Second\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{"--allow-multiple-h1", path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)
	require.Empty(t, stdout.String())
}

func TestRunSkippedLevel(t *testing.T) {
	path := writeTempMarkdown(t, "# Title\n\n## Section\n\n#### Deep\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 1, code)
	require.Contains(t, stdout.String(), "[skipped-level]")
}

func TestRunDuplicateHeading(t *testing.T) {
	path := writeTempMarkdown(t, "# Title\n\n## Section\n\n## Section\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 1, code)
	require.Contains(t, stdout.String(), "[duplicate-heading]")
}

func TestRunNoDuplicates(t *testing.T) {
	path := writeTempMarkdown(t, "# Title\n\n## Section\n\n## Section\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{"--no-duplicates", path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)
	require.Empty(t, stdout.String())
}

func TestRunEmptyHeading(t *testing.T) {
	path := writeTempMarkdown(t, "# Title\n\n## \n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 1, code)
	require.Contains(t, stdout.String(), "[empty-heading]")
}

func TestRunMaxLevel(t *testing.T) {
	path := writeTempMarkdown(t, "# Title\n\n## Section\n\n#### Deep\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{"--max-level", "3", path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 1, code)
	require.Contains(t, stdout.String(), "[max-level]")
}

func TestRunJSONOutput(t *testing.T) {
	path := writeTempMarkdown(t, "## Section\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{"--json", path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 1, code)
	require.Contains(t, stdout.String(), `"rule": "missing-h1"`)
}

func TestRunMissingFile(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{"does-not-exist.md"}, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 2, code)
}

func TestRunMissingInputArgument(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run(nil, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 2, code)
	require.Contains(t, err.Error(), "missing input Markdown file")
}

func TestRunInvalidMaxLevel(t *testing.T) {
	path := writeTempMarkdown(t, "# Title\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{"--max-level", "7", path}, &stdout, &stderr)
	require.EqualError(t, err, "max-level must be at most 6")
	require.Equal(t, 2, code)
}

func TestRunOutputOrderByLine(t *testing.T) {
	path := writeTempMarkdown(t, "# First\n\n#### Deep\n\n# Second\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 1, code)

	output := stdout.String()
	skipped := strings.Index(output, "[skipped-level]")
	multiple := strings.Index(output, "[multiple-h1]")
	require.Greater(t, multiple, skipped)
}

func TestRunOutputOrderByRule(t *testing.T) {
	path := writeTempMarkdown(t, "# First\n\n#### Deep\n\n# Second\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{"--group-by-rule", path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 1, code)

	output := stdout.String()
	skipped := strings.Index(output, "[skipped-level]")
	multiple := strings.Index(output, "[multiple-h1]")
	require.Greater(t, skipped, multiple)
}

func writeTempMarkdown(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "document.md")
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	return path
}
