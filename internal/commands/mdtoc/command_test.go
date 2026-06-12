package mdtoc

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunSimpleDocument(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)

	expected := strings.TrimSpace(`
- [Simple Document](#simple-document)
  - [Getting Started](#getting-started)
    - [Installation](#installation)
    - [Usage](#usage)
  - [Configuration](#configuration)
    - [Output](#output)
  - [Conclusion](#conclusion)
	`)

	require.Equal(t, expected, strings.TrimSpace(stdout.String()))
}

func TestRunCodeBlocksDocument(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "codeblocks.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)

	output := stdout.String()
	require.Contains(t, output, "- [Code Block Test](#code-block-test)")
	require.Contains(t, output, "- [Real Heading](#real-heading)")
	require.NotContains(t, output, "# Fake Heading")
}

func TestRunFrontMatterDocument(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "frontmatter.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{path}, &stdout, &stderr)
	require.NoError(t, err)

	expected := strings.TrimSpace(`
- [Front Matter Test](#front-matter-test)
  - [Section One](#section-one)
  - [Section Two](#section-two)
    - [Details](#details)
  - [Final Notes](#final-notes)
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

func TestRunMinLevel(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--min-level", "2", path}, &stdout, &stderr)
	require.NoError(t, err)

	output := stdout.String()
	require.NotContains(t, output, "Simple Document")
	require.Contains(t, output, "- [Getting Started](#getting-started)")
}

func TestRunMaxLevel(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--max-level", "2", path}, &stdout, &stderr)
	require.NoError(t, err)

	output := stdout.String()
	require.Contains(t, output, "- [Simple Document](#simple-document)")
	require.Contains(t, output, "- [Getting Started](#getting-started)")
	require.NotContains(t, output, "Installation")
}

func TestRunOrderedAndNoLinks(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--ordered", "--no-links", path}, &stdout, &stderr)
	require.NoError(t, err)

	output := stdout.String()
	require.Contains(t, output, "1. Simple Document")
	require.Contains(t, output, "  1.1. Getting Started")
	require.NotContains(t, output, "[")
}

func TestRunNoIndent(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--no-indent", path}, &stdout, &stderr)
	require.NoError(t, err)

	output := stdout.String()
	require.NotContains(t, output, "  -")
	require.Contains(t, output, "- [Simple Document](#simple-document)")
	require.Contains(t, output, "- [Getting Started](#getting-started)")
}

func TestRunNoIndentOrdered(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--no-indent", "--ordered", path}, &stdout, &stderr)
	require.NoError(t, err)

	output := stdout.String()
	require.NotContains(t, output, "  1.")
	require.Contains(t, output, "1. [Simple Document](#simple-document)")
	require.Contains(t, output, "1.1. [Getting Started](#getting-started)")
}

func TestValidateLevels(t *testing.T) {
	tests := []struct {
		name    string
		min     int
		max     int
		wantErr string
	}{
		{
			name: "valid defaults",
			min:  1,
			max:  6,
		},
		{
			name:    "min too low",
			min:     0,
			max:     6,
			wantErr: "min-level must be at least 1",
		},
		{
			name:    "max too high",
			min:     1,
			max:     7,
			wantErr: "max-level must be at most 6",
		},
		{
			name:    "min greater than max",
			min:     4,
			max:     2,
			wantErr: "min-level must not be greater than max-level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateLevels(tt.min, tt.max)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, tt.wantErr)
		})
	}
}

func TestRunInvalidMinLevel(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--min-level", "0", path}, &stdout, &stderr)
	require.EqualError(t, err, "min-level must be at least 1")
}

func TestRunInvalidMaxLevel(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--max-level", "7", path}, &stdout, &stderr)
	require.EqualError(t, err, "max-level must be at most 6")
}

func TestRunInvalidLevelRange(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "simple.md")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run([]string{"--min-level", "4", "--max-level", "2", path}, &stdout, &stderr)
	require.EqualError(t, err, "min-level must not be greater than max-level")
}
