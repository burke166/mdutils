package mdsplit

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeTempMarkdown(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "input.md")
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))
	return path
}

func TestRunSplitsByH1Default(t *testing.T) {
	content := "# Introduction\n\nWelcome.\n\n# Character Creation\n\nMake a character.\n\n## Abilities\n\nAbility rules.\n\n# Combat\n\nCombat rules.\n"
	path := writeTempMarkdown(t, content)
	outDir := filepath.Join(t.TempDir(), "chapters")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--out", outDir, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	intro, err := os.ReadFile(filepath.Join(outDir, "introduction.md"))
	require.NoError(t, err)
	require.Equal(t, "# Introduction\n\nWelcome.\n", string(intro))

	creation, err := os.ReadFile(filepath.Join(outDir, "character-creation.md"))
	require.NoError(t, err)
	require.Equal(t, "# Character Creation\n\nMake a character.\n\n## Abilities\n\nAbility rules.\n", string(creation))

	combat, err := os.ReadFile(filepath.Join(outDir, "combat.md"))
	require.NoError(t, err)
	require.Equal(t, "# Combat\n\nCombat rules.\n", string(combat))
}

func TestRunSplitsByH2(t *testing.T) {
	content := "# Book\n\nIntro.\n\n## Chapter One\n\nFirst.\n\n## Chapter Two\n\nSecond.\n"
	path := writeTempMarkdown(t, content)
	outDir := t.TempDir()

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--level", "2", "--out", outDir, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	one, err := os.ReadFile(filepath.Join(outDir, "chapter-one.md"))
	require.NoError(t, err)
	require.Equal(t, "## Chapter One\n\nFirst.\n", string(one))
}

func TestRunWritesPreamble(t *testing.T) {
	content := "Preamble.\n\n# First\n\nSection.\n"
	path := writeTempMarkdown(t, content)
	outDir := t.TempDir()

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--out", outDir, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	preamble, err := os.ReadFile(filepath.Join(outDir, "00-preamble.md"))
	require.NoError(t, err)
	require.Equal(t, "Preamble.\n", string(preamble))
}

func TestRunCreatesOutputDirectory(t *testing.T) {
	content := "# Only\n\nBody.\n"
	path := writeTempMarkdown(t, content)
	outDir := filepath.Join(t.TempDir(), "nested", "chapters")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--out", outDir, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	_, err = os.Stat(filepath.Join(outDir, "only.md"))
	require.NoError(t, err)
}

func TestRunSanitizesIllegalFilenameCharacters(t *testing.T) {
	content := "# foo:bar/baz\n\nBody.\n"
	path := writeTempMarkdown(t, content)
	outDir := t.TempDir()

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--out", outDir, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	_, err = os.Stat(filepath.Join(outDir, "foo-bar-baz.md"))
	require.NoError(t, err)
}

func TestRunHandlesDuplicateHeadingSlugs(t *testing.T) {
	content := "# Introduction\n\nFirst.\n\n# Introduction\n\nSecond.\n"
	path := writeTempMarkdown(t, content)
	outDir := t.TempDir()

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--out", outDir, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	_, err = os.Stat(filepath.Join(outDir, "introduction.md"))
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(outDir, "introduction-2.md"))
	require.NoError(t, err)
}

func TestRunNoMatchingHeadings(t *testing.T) {
	path := writeTempMarkdown(t, "## Only H2\n\nBody.\n")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{path}, &stdout, &stderr)
	require.Error(t, err)
	require.Equal(t, 1, code)
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

func TestRunIgnoresHeadingsInCodeBlocks(t *testing.T) {
	content := "# Real\n\n```\n# Fake\n```\n\nBody.\n"
	path := writeTempMarkdown(t, content)
	outDir := t.TempDir()

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--out", outDir, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	data, err := os.ReadFile(filepath.Join(outDir, "real.md"))
	require.NoError(t, err)
	require.Contains(t, string(data), "# Fake")
}

func TestRunDefaultFilenames(t *testing.T) {
	content := "# Intro\n\nWelcome.\n\n# Usage\n\nHow to use.\n\n# FAQ\n\nQuestions.\n"
	path := writeTempMarkdown(t, content)
	outDir := t.TempDir()

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--out", outDir, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	for _, name := range []string{"intro.md", "usage.md", "faq.md"} {
		_, err := os.Stat(filepath.Join(outDir, name))
		require.NoError(t, err, "expected %s", name)
	}
}

func TestRunNumberedFilenames(t *testing.T) {
	content := "# Intro\n\nWelcome.\n\n# Usage\n\nHow to use.\n\n# FAQ\n\nQuestions.\n"
	path := writeTempMarkdown(t, content)
	outDir := t.TempDir()

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--numbered", "--out", outDir, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	for _, name := range []string{"01_intro.md", "02_usage.md", "03_faq.md"} {
		_, err := os.Stat(filepath.Join(outDir, name))
		require.NoError(t, err, "expected %s", name)
	}
}

func TestRunNumberedShortFlag(t *testing.T) {
	content := "# Intro\n\nWelcome.\n\n# Usage\n\nHow to use.\n"
	path := writeTempMarkdown(t, content)
	outDir := t.TempDir()

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"-n", "--out", outDir, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	_, err = os.Stat(filepath.Join(outDir, "01_intro.md"))
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(outDir, "02_usage.md"))
	require.NoError(t, err)
}

func TestRunNumberedWidthOver99(t *testing.T) {
	var content strings.Builder
	for i := 1; i <= 100; i++ {
		content.WriteString("# File\n\nBody.\n\n")
	}
	path := writeTempMarkdown(t, content.String())
	outDir := t.TempDir()

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--numbered", "--out", outDir, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	_, err = os.Stat(filepath.Join(outDir, "001_file.md"))
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(outDir, "100_file-100.md"))
	require.NoError(t, err)
}

func TestRunNumberedSanitizesSlug(t *testing.T) {
	content := "# Chapter: One?\n\nBody.\n"
	path := writeTempMarkdown(t, content)
	outDir := t.TempDir()

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--numbered", "--out", outDir, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	_, err = os.Stat(filepath.Join(outDir, "01_chapter-one.md"))
	require.NoError(t, err)
}

func TestRunNumberedOutputDirectory(t *testing.T) {
	content := "# Intro\n\nWelcome.\n\n# Usage\n\nHow to use.\n"
	path := writeTempMarkdown(t, content)
	outDir := filepath.Join(t.TempDir(), "docs", "chapters")

	var stdout, stderr bytes.Buffer
	code, err := Run([]string{"--numbered", "--out", outDir, path}, &stdout, &stderr)
	require.NoError(t, err)
	require.Equal(t, 0, code)

	_, err = os.Stat(filepath.Join(outDir, "01_intro.md"))
	require.NoError(t, err)
	_, err = os.Stat(filepath.Join(outDir, "02_usage.md"))
	require.NoError(t, err)
}
