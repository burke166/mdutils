package mdformat

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeTempMarkdown(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "input.md")
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}

func TestRunPrintsFormattedOutputToStdout(t *testing.T) {
	path := writeTempMarkdown(t, "##Heading\n")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)
	require.Equal(t, "## Heading\n", stdout.String())
}

func TestRunCheckExitsZeroWhenFormatted(t *testing.T) {
	path := writeTempMarkdown(t, "# Title\n\nBody\n")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--check", path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)
	require.Empty(t, stdout.String())
}

func TestRunCheckExitsOneWhenChangesNeeded(t *testing.T) {
	path := writeTempMarkdown(t, "##Heading\n")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--check", path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 1, code)
	require.Empty(t, stdout.String())
}

func TestRunWriteUpdatesFile(t *testing.T) {
	path := writeTempMarkdown(t, "##Heading\n")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--write", path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "## Heading\n", string(data))
}

func TestRunOutputWritesSeparateFile(t *testing.T) {
	path := writeTempMarkdown(t, "##Heading\n")
	outPath := filepath.Join(t.TempDir(), "formatted.md")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--output", outPath, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	data, err := os.ReadFile(outPath)
	require.NoError(t, err)
	require.Equal(t, "## Heading\n", string(data))

	original, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "##Heading\n", string(original))
}

func TestRunWritePreservesPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permission bits are not preserved on Windows")
	}
	path := writeTempMarkdown(t, "##Heading\n")
	require.NoError(t, os.Chmod(path, 0600))

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--write", path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	info, err := os.Stat(path)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0600), info.Mode().Perm())

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, "## Heading\n", string(data))
}

func TestRunRejectsWriteAndOutputTogether(t *testing.T) {
	path := writeTempMarkdown(t, "# Title\n")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--write", "--output", "out.md", path}, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 2, code)
	require.Contains(t, err.Error(), "choose only one")
}

func TestRunRejectsCheckAndWriteTogether(t *testing.T) {
	path := writeTempMarkdown(t, "# Title\n")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--check", "--write", path}, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 2, code)
	require.Contains(t, err.Error(), "choose only one")
}

func TestRunRejectsOutputSameAsInput(t *testing.T) {
	path := writeTempMarkdown(t, "##Heading\n")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--output", path, path}, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 2, code)
	require.Contains(t, err.Error(), "must differ")
}

func TestRunMissingInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code, err := Run(nil, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 2, code)
	require.Contains(t, err.Error(), "missing input Markdown file")
}

func TestRunMissingFile(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"missing.md"}, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 2, code)
}
