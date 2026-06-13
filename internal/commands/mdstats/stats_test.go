package mdstats

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func testdataPath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join("..", "..", "..", "testdata", name)
}

func TestAnalyzeSimpleDocument(t *testing.T) {
	path := testdataPath(t, "simple.md")
	content, err := os.ReadFile(path)
	require.NoError(t, err)

	stats := AnalyzeContent(path, content)
	require.Equal(t, path, stats.Path)
	require.Greater(t, stats.Lines, 0)
	require.Greater(t, stats.Words, 0)
	require.Equal(t, 1, stats.Headings.H1)
	require.Equal(t, 3, stats.Headings.H2)
	require.Equal(t, 3, stats.Headings.H3)
	require.Equal(t, 7, stats.Headings.Total)
	require.Equal(t, 3, stats.Headings.MaxDepth)
	require.Equal(t, 3, stats.Markdown.NumberedItems)
	require.Equal(t, 1, stats.Markdown.CodeBlocks)
	require.False(t, stats.Frontmatter.Detected)
}

func TestAnalyzeFrontmatterDetection(t *testing.T) {
	path := testdataPath(t, "frontmatter.md")
	content, err := os.ReadFile(path)
	require.NoError(t, err)

	stats := AnalyzeContent(path, content)
	require.True(t, stats.Frontmatter.Detected)
	require.Equal(t, 9, stats.Frontmatter.Lines)
	require.Equal(t, 1, stats.Headings.H1)
	require.Equal(t, 3, stats.Headings.H2)
	require.Equal(t, 1, stats.Headings.H3)
	require.Equal(t, 0, stats.Headings.H4)
}

func TestAnalyzeCodeFenceExclusion(t *testing.T) {
	path := testdataPath(t, "codeblocks.md")
	content, err := os.ReadFile(path)
	require.NoError(t, err)

	stats := AnalyzeContent(path, content)
	require.Equal(t, 1, stats.Headings.H1)
	require.Equal(t, 4, stats.Headings.H2)
	require.Equal(t, 0, stats.Headings.H3)
	require.Equal(t, 5, stats.Headings.Total)
	require.Equal(t, 3, stats.Markdown.CodeBlocks)
}

func TestAnalyzeHeadingLevels(t *testing.T) {
	content := []byte(`# One
## Two
### Three
#### Four
##### Five
###### Six
####### Not a heading
#No space
`)

	stats := AnalyzeContent("headings.md", content)
	require.Equal(t, 1, stats.Headings.H1)
	require.Equal(t, 1, stats.Headings.H2)
	require.Equal(t, 1, stats.Headings.H3)
	require.Equal(t, 1, stats.Headings.H4)
	require.Equal(t, 1, stats.Headings.H5)
	require.Equal(t, 1, stats.Headings.H6)
	require.Equal(t, 6, stats.Headings.Total)
	require.Equal(t, 6, stats.Headings.MaxDepth)
}

func TestAnalyzeLinksAndImages(t *testing.T) {
	content := []byte(`[link](https://example.com)
![image](pic.png)
[ref][id]
![ref image][id]
`)

	stats := AnalyzeContent("links.md", content)
	require.Equal(t, 2, stats.Markdown.Links)
	require.Equal(t, 2, stats.Markdown.Images)
}

func TestAnalyzeListsTaskListsBlockquotesTablesFootnotesRules(t *testing.T) {
	content := []byte(`- bullet
* another
+ third
1. numbered
2. second
- [ ] task open
- [x] task done

> quote one
> quote two

| Col A | Col B |
| ----- | ----- |
| a     | b     |

See footnote[^note].

---

***

___
`)

	stats := AnalyzeContent("features.md", content)
	require.Equal(t, 7, stats.Markdown.Lists)
	require.Equal(t, 5, stats.Markdown.BulletItems)
	require.Equal(t, 2, stats.Markdown.NumberedItems)
	require.Equal(t, 2, stats.Markdown.TaskItems)
	require.Equal(t, 2, stats.Markdown.BlockQuoteLines)
	require.Equal(t, 1, stats.Markdown.Tables)
	require.Equal(t, 1, stats.Markdown.Footnotes)
	require.Equal(t, 3, stats.Markdown.HorizontalRules)
}

func TestAnalyzeInlineCodeSpans(t *testing.T) {
	content := []byte("Use `code` and ``two ticks`` here.\n")
	stats := AnalyzeContent("inline.md", content)
	require.Equal(t, 2, stats.Markdown.InlineCodeSpans)
}

func TestReadingTimeCalculation(t *testing.T) {
	require.Equal(t, 0, readingTimeMinutes(0))
	require.Equal(t, 1, readingTimeMinutes(1))
	require.Equal(t, 1, readingTimeMinutes(250))
	require.Equal(t, 2, readingTimeMinutes(251))
	require.Equal(t, 14, readingTimeMinutes(3412))
}

func TestSummarizeAggregation(t *testing.T) {
	files := []FileStats{
		{Words: 100, Lines: 10, Sentences: 4, Paragraphs: 2, Characters: 50, BlankLines: 1, FileSizeBytes: 100, ReadingTimeMinutes: 1},
		{Words: 200, Lines: 20, Sentences: 8, Paragraphs: 3, Characters: 80, BlankLines: 2, FileSizeBytes: 200, ReadingTimeMinutes: 1},
	}

	summary := Summarize(files)
	require.Equal(t, 2, summary.FileCount)
	require.Equal(t, 300, summary.Words)
	require.Equal(t, 30, summary.Lines)
	require.Equal(t, 12, summary.Sentences)
	require.Equal(t, 5, summary.Paragraphs)
	require.Equal(t, 130, summary.Characters)
	require.Equal(t, 3, summary.BlankLines)
	require.Equal(t, int64(300), summary.FileSizeBytes)
	require.Equal(t, 2, summary.ReadingTimeMinutes)
}

func TestCountSentences(t *testing.T) {
	require.Equal(t, 0, countSentences(""))
	require.Equal(t, 1, countSentences("Hello world"))
	require.Equal(t, 2, countSentences("Hello. World!"))
}

func TestCountWordsStripsMarkdown(t *testing.T) {
	require.Equal(t, 2, countWords("[link](https://example.com) text"))
	require.Equal(t, 2, countWords("![image](pic.png) words"))
}
