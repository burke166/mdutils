package mdsplit

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitMarkdownByH1(t *testing.T) {
	content := "# Introduction\n\nWelcome.\n\n# Character Creation\n\nMake a character.\n\n# Combat\n\nCombat rules.\n"

	sections, err := SplitMarkdown(content, 1)
	require.NoError(t, err)
	require.Len(t, sections, 3)
	require.Equal(t, "Introduction", sections[0].Heading)
	require.Equal(t, "introduction", sections[0].Slug)
	require.Equal(t, "# Introduction\n\nWelcome.\n", sections[0].Content)
	require.Equal(t, "Character Creation", sections[1].Heading)
	require.Equal(t, "# Character Creation\n\nMake a character.\n", sections[1].Content)
	require.Equal(t, "Combat", sections[2].Heading)
}

func TestSplitMarkdownByH2(t *testing.T) {
	content := "## Chapter One\n\nFirst chapter.\n\n## Chapter Two\n\nSecond chapter.\n"

	sections, err := SplitMarkdown(content, 2)
	require.NoError(t, err)
	require.Len(t, sections, 2)
	require.Equal(t, "Chapter One", sections[0].Heading)
	require.Equal(t, "## Chapter One\n\nFirst chapter.\n", sections[0].Content)
	require.Equal(t, "Chapter Two", sections[1].Heading)
}

func TestSplitMarkdownPreservesNestedHeadings(t *testing.T) {
	content := "# Character Creation\n\nMake a character.\n\n## Abilities\n\nAbility rules.\n"

	sections, err := SplitMarkdown(content, 1)
	require.NoError(t, err)
	require.Len(t, sections, 1)
	require.Equal(t, content, sections[0].Content)
}

func TestSplitMarkdownIgnoresHeadingsInCodeBlocks(t *testing.T) {
	content := "# Real Heading\n\n```\n# not a heading\n```\n\nBody.\n"

	sections, err := SplitMarkdown(content, 1)
	require.NoError(t, err)
	require.Len(t, sections, 1)
	require.Contains(t, sections[0].Content, "# not a heading")
}

func TestSplitMarkdownPreamble(t *testing.T) {
	content := "Preamble text.\n\n# First\n\nSection one.\n"

	sections, err := SplitMarkdown(content, 1)
	require.NoError(t, err)
	require.Len(t, sections, 2)
	require.Empty(t, sections[0].Heading)
	require.Equal(t, "Preamble text.\n", sections[0].Content)
	require.Equal(t, "First", sections[1].Heading)
}

func TestSplitMarkdownNoMatchingHeadings(t *testing.T) {
	_, err := SplitMarkdown("## Only H2\n\nBody.\n", 1)
	require.EqualError(t, err, "no matching headings found")
}

func TestSlugifyFilename(t *testing.T) {
	tests := []struct {
		name     string
		heading  string
		expected string
	}{
		{name: "simple", heading: "Introduction", expected: "introduction"},
		{name: "punctuation", heading: "Character Creation!", expected: "character-creation"},
		{name: "illegal chars", heading: `foo:bar/baz\qux|?*<>""`, expected: "foo-bar-baz-qux"},
		{name: "collapse dashes", heading: "foo   bar", expected: "foo-bar"},
		{name: "empty", heading: "!!!", expected: "section"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, SlugifyFilename(tt.heading))
		})
	}
}

func TestEnsureUniqueFilename(t *testing.T) {
	used := make(map[string]int)

	require.Equal(t, "introduction", EnsureUniqueFilename("introduction", used))
	require.Equal(t, "introduction-2", EnsureUniqueFilename("introduction", used))
	require.Equal(t, "introduction-3", EnsureUniqueFilename("introduction", used))
}

func TestSplitMarkdownPreservesBlankLines(t *testing.T) {
	content := "# Title\n\n\n\nParagraph with blank lines.\n\n# Next\n\nDone.\n"

	sections, err := SplitMarkdown(content, 1)
	require.NoError(t, err)
	require.Len(t, sections, 2)
	require.Contains(t, sections[0].Content, "\n\n\n")
}
