package mdlint

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeMarkdown(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}

func TestRunValidSingleFile(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)
	require.Empty(t, stdout.String())
}

func TestRunSingleFileWithIssues(t *testing.T) {
	dir := t.TempDir()
	path := writeMarkdown(t, dir, "doc.md", "# First\n\n# Second\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 2, code)
	require.Contains(t, stdout.String(), "single-h1")
}

func TestRunFolderRecursive(t *testing.T) {
	dir := t.TempDir()
	writeMarkdown(t, dir, "a.md", "# A\n\nContent.\n")
	nested := filepath.Join(dir, "nested")
	require.NoError(t, os.Mkdir(nested, 0755))
	writeMarkdown(t, nested, "b.md", "## Missing H1\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{dir}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 1, code)
	require.Contains(t, stdout.String(), "no-missing-h1")
}

func TestRunFolderNoRecursive(t *testing.T) {
	dir := t.TempDir()
	writeMarkdown(t, dir, "a.md", "# A\n\nContent.\n")
	nested := filepath.Join(dir, "nested")
	require.NoError(t, os.Mkdir(nested, 0755))
	writeMarkdown(t, nested, "b.md", "## Missing H1\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{"--no-recursive", dir}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)
	require.Empty(t, stdout.String())
}

func TestRunIgnoredFolders(t *testing.T) {
	dir := t.TempDir()
	writeMarkdown(t, dir, "a.md", "# A\n\nContent.\n")
	vendor := filepath.Join(dir, "vendor")
	require.NoError(t, os.Mkdir(vendor, 0755))
	writeMarkdown(t, vendor, "b.md", "## Missing H1\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{dir}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)
}

func TestRunExcludeGlobsFromConfig(t *testing.T) {
	dir := t.TempDir()
	writeMarkdown(t, dir, "README.md", "# Title\n\nContent.\n")
	writeMarkdown(t, dir, "CHANGELOG.md", "## Missing H1\n")

	configPath := filepath.Join(dir, "config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
exclude:
  - "CHANGELOG.md"
`), 0644))

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{"--config", configPath, dir}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)
}

func TestRunJSONOutput(t *testing.T) {
	dir := t.TempDir()
	path := writeMarkdown(t, dir, "doc.md", "## Missing H1\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{"--json", path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 1, code)
	require.Contains(t, stdout.String(), `"ruleId": "no-missing-h1"`)
}

func TestRunQuietOutput(t *testing.T) {
	dir := t.TempDir()
	path := writeMarkdown(t, dir, "doc.md", "## Missing H1\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{"--quiet", path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 1, code)
	require.Empty(t, stdout.String())
}

func TestRunWarningsOnlyExitCode(t *testing.T) {
	dir := t.TempDir()
	path := writeMarkdown(t, dir, "doc.md", "## Missing H1\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 1, code)
}

func TestRunErrorsExitCode(t *testing.T) {
	dir := t.TempDir()
	path := writeMarkdown(t, dir, "doc.md", "# First\n\n# Second\n")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 2, code)
}

func TestRunMissingInput(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run(nil, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 3, code)
}

func TestRunMissingFile(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{"missing.md"}, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 3, code)
}

func TestRunInvalidConfig(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "bad.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte("rules: ["), 0644))

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{"--config", configPath, dir}, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 3, code)
}

func TestRunUsesExplicitConfigOnly(t *testing.T) {
	dir := t.TempDir()
	writeMarkdown(t, dir, "doc.md", "# First\n\nIntro.\n\n# Second\n\nMore.\n")

	localConfig := filepath.Join(dir, configFileName)
	require.NoError(t, os.WriteFile(localConfig, []byte(`
rules:
  single-h1: true
`), 0644))

	otherDir := t.TempDir()
	otherConfig := filepath.Join(otherDir, "other.yaml")
	require.NoError(t, os.WriteFile(otherConfig, []byte(`
rules:
  single-h1: false
`), 0644))

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{"--config", otherConfig, filepath.Join(dir, "doc.md")}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)
}

func TestRunAllFilesExcluded(t *testing.T) {
	dir := t.TempDir()
	writeMarkdown(t, dir, "CHANGELOG.md", "## Missing H1\n")

	configPath := filepath.Join(dir, "config.yaml")
	require.NoError(t, os.WriteFile(configPath, []byte(`
exclude:
  - "CHANGELOG.md"
`), 0644))

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code, err := Run([]string{"--config", configPath, dir}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)
}
