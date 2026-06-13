package mdstats

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func RenderText(result AnalysisResult, perFile bool, summary bool) string {
	var b strings.Builder

	if perFile {
		for i, file := range result.Files {
			if i > 0 {
				b.WriteString("\n")
			}
			renderFileText(&b, file)
		}
	}

	if summary {
		if perFile && len(result.Files) > 0 {
			b.WriteString("\n")
		}
		renderSummaryText(&b, result.Summary)
	}

	return b.String()
}

func renderFileText(b *strings.Builder, file FileStats) {
	fmt.Fprintf(b, "%s\n\n", file.Path)

	writeLabelInt(b, "Lines:", file.Lines)
	writeLabelInt(b, "Blank lines:", file.BlankLines)
	writeLabelInt(b, "Words:", file.Words)
	writeLabelInt(b, "Characters:", file.Characters)
	fmt.Fprintf(b, "File size:          %s bytes\n", formatInt64(file.FileSizeBytes))
	writeLabelInt(b, "Paragraphs:", file.Paragraphs)
	writeLabelInt(b, "Sentences:", file.Sentences)
	fmt.Fprintf(b, "Reading time:       %s\n", formatMinutes(file.ReadingTimeMinutes))

	b.WriteString("\nHeadings:\n")
	writeLabelInt(b, "  H1:", file.Headings.H1)
	writeLabelInt(b, "  H2:", file.Headings.H2)
	writeLabelInt(b, "  H3:", file.Headings.H3)
	writeLabelInt(b, "  H4:", file.Headings.H4)
	writeLabelInt(b, "  H5:", file.Headings.H5)
	writeLabelInt(b, "  H6:", file.Headings.H6)

	b.WriteString("\nMarkdown:\n")
	writeLabelInt(b, "  Lists:", file.Markdown.Lists)
	writeLabelInt(b, "  Bullet items:", file.Markdown.BulletItems)
	writeLabelInt(b, "  Numbered items:", file.Markdown.NumberedItems)
	writeLabelInt(b, "  Task items:", file.Markdown.TaskItems)
	writeLabelInt(b, "  Block quote lines:", file.Markdown.BlockQuoteLines)
	writeLabelInt(b, "  Code blocks:", file.Markdown.CodeBlocks)
	writeLabelInt(b, "  Inline code spans:", file.Markdown.InlineCodeSpans)
	writeLabelInt(b, "  Tables:", file.Markdown.Tables)
	writeLabelInt(b, "  Links:", file.Markdown.Links)
	writeLabelInt(b, "  Images:", file.Markdown.Images)
	writeLabelInt(b, "  Footnotes:", file.Markdown.Footnotes)
	writeLabelInt(b, "  Horizontal rules:", file.Markdown.HorizontalRules)

	b.WriteString("\nFrontmatter:\n")
	fmt.Fprintf(b, "  Detected: %t\n", file.Frontmatter.Detected)
	writeLabelInt(b, "  Lines:", file.Frontmatter.Lines)
}

func renderSummaryText(b *strings.Builder, summary SummaryStats) {
	b.WriteString("Summary\n\n")
	writeLabelInt(b, "Files:", summary.FileCount)
	writeLabelInt(b, "Lines:", summary.Lines)
	writeLabelInt(b, "Blank lines:", summary.BlankLines)
	writeLabelInt(b, "Words:", summary.Words)
	writeLabelInt(b, "Characters:", summary.Characters)
	fmt.Fprintf(b, "File size:          %s bytes\n", formatInt64(summary.FileSizeBytes))
	writeLabelInt(b, "Paragraphs:", summary.Paragraphs)
	writeLabelInt(b, "Sentences:", summary.Sentences)
	fmt.Fprintf(b, "Reading time:       %s\n", formatMinutes(summary.ReadingTimeMinutes))
}

func RenderJSON(result AnalysisResult) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

func RenderCSV(files []FileStats) (string, error) {
	var b strings.Builder
	w := csv.NewWriter(&b)

	header := []string{
		"path", "fileSizeBytes", "lines", "blankLines", "characters", "words",
		"paragraphs", "sentences", "readingTimeMinutes",
		"h1", "h2", "h3", "h4", "h5", "h6", "totalHeadings", "maxHeadingDepth",
		"lists", "bulletItems", "numberedItems", "taskItems", "blockQuoteLines",
		"codeBlocks", "inlineCodeSpans", "tables", "links", "images", "footnotes",
		"horizontalRules", "frontmatterDetected", "frontmatterLines",
	}
	if err := w.Write(header); err != nil {
		return "", err
	}

	for _, file := range files {
		row := []string{
			file.Path,
			strconv.FormatInt(file.FileSizeBytes, 10),
			strconv.Itoa(file.Lines),
			strconv.Itoa(file.BlankLines),
			strconv.Itoa(file.Characters),
			strconv.Itoa(file.Words),
			strconv.Itoa(file.Paragraphs),
			strconv.Itoa(file.Sentences),
			strconv.Itoa(file.ReadingTimeMinutes),
			strconv.Itoa(file.Headings.H1),
			strconv.Itoa(file.Headings.H2),
			strconv.Itoa(file.Headings.H3),
			strconv.Itoa(file.Headings.H4),
			strconv.Itoa(file.Headings.H5),
			strconv.Itoa(file.Headings.H6),
			strconv.Itoa(file.Headings.Total),
			strconv.Itoa(file.Headings.MaxDepth),
			strconv.Itoa(file.Markdown.Lists),
			strconv.Itoa(file.Markdown.BulletItems),
			strconv.Itoa(file.Markdown.NumberedItems),
			strconv.Itoa(file.Markdown.TaskItems),
			strconv.Itoa(file.Markdown.BlockQuoteLines),
			strconv.Itoa(file.Markdown.CodeBlocks),
			strconv.Itoa(file.Markdown.InlineCodeSpans),
			strconv.Itoa(file.Markdown.Tables),
			strconv.Itoa(file.Markdown.Links),
			strconv.Itoa(file.Markdown.Images),
			strconv.Itoa(file.Markdown.Footnotes),
			strconv.Itoa(file.Markdown.HorizontalRules),
			strconv.FormatBool(file.Frontmatter.Detected),
			strconv.Itoa(file.Frontmatter.Lines),
		}
		if err := w.Write(row); err != nil {
			return "", err
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}

	return b.String(), nil
}

func writeLabelInt(b *strings.Builder, label string, value int) {
	fmt.Fprintf(b, "%-20s %s\n", label, formatInt(value))
}

func formatInt(n int) string {
	s := strconv.Itoa(n)
	if n < 1000 {
		return s
	}

	var parts []string
	for s != "" {
		start := len(s) - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{s[start:]}, parts...)
		s = s[:start]
	}
	return strings.Join(parts, ",")
}

func formatInt64(n int64) string {
	return formatInt(int(n))
}

func formatMinutes(minutes int) string {
	if minutes == 1 {
		return "1 minute"
	}
	return fmt.Sprintf("%s minutes", formatInt(minutes))
}
