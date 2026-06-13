package mdstats

import (
	"math"
	"os"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

const wordsPerMinute = 250

var (
	markdownLinkPattern    = regexp.MustCompile(`\[([^\]]*)\]\(([^)]*)\)`)
	markdownLinkRefPattern = regexp.MustCompile(`\[([^\]]*)\]\[([^\]]*)\]`)
	imagePattern           = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]*)\)`)
	imageRefPattern        = regexp.MustCompile(`!\[([^\]]*)\]\[([^\]]*)\]`)
	footnoteRefPattern     = regexp.MustCompile(`\[\^[^\]]+\]`)
	taskListPattern        = regexp.MustCompile(`^\s*[-*+]\s+\[[ xX]\]\s+`)
	bulletListPattern      = regexp.MustCompile(`^\s*[-*+]\s+`)
	numberedListPattern    = regexp.MustCompile(`^\s*\d+\.\s+`)
	blockQuotePattern      = regexp.MustCompile(`^\s*>`)
	tableSeparatorPattern  = regexp.MustCompile(`^\s*\|?(\s*:?-{3,}:?\s*\|)+\s*:?-{3,}:?\s*\|?\s*$`)
)

type FileStats struct {
	Path               string           `json:"path"`
	FileSizeBytes      int64            `json:"fileSizeBytes"`
	Lines              int              `json:"lines"`
	BlankLines         int              `json:"blankLines"`
	Characters         int              `json:"characters"`
	Words              int              `json:"words"`
	Paragraphs         int              `json:"paragraphs"`
	Sentences          int              `json:"sentences"`
	ReadingTimeMinutes int              `json:"readingTimeMinutes"`
	Headings           HeadingStats     `json:"headings"`
	Markdown           MarkdownStats    `json:"markdown"`
	Frontmatter        FrontmatterStats `json:"frontmatter"`
}

type HeadingStats struct {
	H1       int `json:"h1"`
	H2       int `json:"h2"`
	H3       int `json:"h3"`
	H4       int `json:"h4"`
	H5       int `json:"h5"`
	H6       int `json:"h6"`
	Total    int `json:"total"`
	MaxDepth int `json:"maxDepth"`
}

type MarkdownStats struct {
	Lists           int `json:"lists"`
	BulletItems     int `json:"bulletItems"`
	NumberedItems   int `json:"numberedItems"`
	TaskItems       int `json:"taskItems"`
	BlockQuoteLines int `json:"blockQuoteLines"`
	CodeBlocks      int `json:"codeBlocks"`
	InlineCodeSpans int `json:"inlineCodeSpans"`
	Tables          int `json:"tables"`
	Links           int `json:"links"`
	Images          int `json:"images"`
	Footnotes       int `json:"footnotes"`
	HorizontalRules int `json:"horizontalRules"`
}

type FrontmatterStats struct {
	Detected bool `json:"detected"`
	Lines    int  `json:"lines"`
}

type SummaryStats struct {
	FileCount          int   `json:"fileCount"`
	FileSizeBytes      int64 `json:"fileSizeBytes"`
	Lines              int   `json:"lines"`
	BlankLines         int   `json:"blankLines"`
	Characters         int   `json:"characters"`
	Words              int   `json:"words"`
	Paragraphs         int   `json:"paragraphs"`
	Sentences          int   `json:"sentences"`
	ReadingTimeMinutes int   `json:"readingTimeMinutes"`
}

type AnalysisResult struct {
	Files   []FileStats  `json:"files"`
	Summary SummaryStats `json:"summary"`
}

func AnalyzePath(path string) (FileStats, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileStats{}, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return FileStats{}, err
	}

	stats := AnalyzeContent(path, content)
	stats.FileSizeBytes = info.Size()
	return stats, nil
}

func AnalyzeContent(path string, content []byte) FileStats {
	text := string(content)
	lines := markdown.SplitLines(text)

	stats := FileStats{
		Path:       path,
		Characters: utf8.RuneCountInString(text),
		Lines:      len(lines),
	}

	for _, line := range lines {
		if markdown.IsBlankLine(line) {
			stats.BlankLines++
		}
	}

	frontmatterLines, hasFrontmatter := markdown.FrontmatterBounds(lines)
	stats.Frontmatter.Detected = hasFrontmatter
	stats.Frontmatter.Lines = frontmatterLines

	bodyStart := 0
	if hasFrontmatter {
		bodyStart = frontmatterLines
	}

	var prose strings.Builder
	inFence := false
	var fenceChar byte
	inTable := false
	prevProse := false

	for i, line := range lines {
		if i < bodyStart {
			continue
		}

		body := strings.TrimRight(line, "\r\n")

		isFence, char := markdown.IsFenceLine(line)
		if isFence {
			if !inFence {
				stats.Markdown.CodeBlocks++
				inFence = true
				fenceChar = char
			} else if char == fenceChar {
				inFence = false
				fenceChar = 0
			}
			inTable = false
			prevProse = false
			continue
		}

		if inFence {
			continue
		}

		stats.Markdown.InlineCodeSpans += countInlineCodeSpans(body)
		stats.Markdown.Links += countLinks(body)
		stats.Markdown.Images += countMatches(imagePattern, body) + countMatches(imageRefPattern, body)
		stats.Markdown.Footnotes += countMatches(footnoteRefPattern, body)

		if level, ok := parseATXHeading(body); ok {
			recordHeading(&stats.Headings, level)
			inTable = false
			prevProse = false
			continue
		}

		if isHorizontalRule(body) {
			stats.Markdown.HorizontalRules++
			inTable = false
			prevProse = false
			continue
		}

		if taskListPattern.MatchString(body) {
			stats.Markdown.TaskItems++
			stats.Markdown.BulletItems++
			stats.Markdown.Lists++
			inTable = false
			prevProse = false
			continue
		}

		if bulletListPattern.MatchString(body) {
			stats.Markdown.BulletItems++
			stats.Markdown.Lists++
			inTable = false
			prevProse = false
			continue
		}

		if numberedListPattern.MatchString(body) {
			stats.Markdown.NumberedItems++
			stats.Markdown.Lists++
			inTable = false
			prevProse = false
			continue
		}

		if blockQuotePattern.MatchString(body) {
			stats.Markdown.BlockQuoteLines++
			inTable = false
			prevProse = false
			continue
		}

		if isTableRow(body) {
			if !inTable && i+1 < len(lines) {
				nextBody := strings.TrimRight(lines[i+1], "\r\n")
				if tableSeparatorPattern.MatchString(nextBody) {
					stats.Markdown.Tables++
				}
			}
			inTable = true
			prevProse = false
			continue
		}

		inTable = false

		if markdown.IsBlankLine(line) {
			prevProse = false
			continue
		}

		if !prevProse {
			stats.Paragraphs++
			prevProse = true
		}

		if prose.Len() > 0 {
			prose.WriteByte(' ')
		}
		prose.WriteString(body)
	}

	wordText := prose.String()
	stats.Words = countWords(wordText)
	stats.Sentences = countSentences(wordText)
	stats.ReadingTimeMinutes = readingTimeMinutes(stats.Words)

	return stats
}

func Summarize(files []FileStats) SummaryStats {
	summary := SummaryStats{FileCount: len(files)}
	for _, file := range files {
		summary.FileSizeBytes += file.FileSizeBytes
		summary.Lines += file.Lines
		summary.BlankLines += file.BlankLines
		summary.Characters += file.Characters
		summary.Words += file.Words
		summary.Paragraphs += file.Paragraphs
		summary.Sentences += file.Sentences
	}
	summary.ReadingTimeMinutes = readingTimeMinutes(summary.Words)
	return summary
}

func readingTimeMinutes(words int) int {
	if words == 0 {
		return 0
	}
	return int(math.Ceil(float64(words) / wordsPerMinute))
}

func recordHeading(headings *HeadingStats, level int) {
	switch level {
	case 1:
		headings.H1++
	case 2:
		headings.H2++
	case 3:
		headings.H3++
	case 4:
		headings.H4++
	case 5:
		headings.H5++
	case 6:
		headings.H6++
	}

	headings.Total++
	if level > headings.MaxDepth {
		headings.MaxDepth = level
	}
}

func parseATXHeading(line string) (int, bool) {
	trimmed := strings.TrimLeft(line, " \t")
	if !strings.HasPrefix(trimmed, "#") {
		return 0, false
	}

	hashes := 0
	for hashes < len(trimmed) && trimmed[hashes] == '#' {
		hashes++
	}
	if hashes < 1 || hashes > 6 {
		return 0, false
	}

	if hashes == len(trimmed) {
		return hashes, true
	}

	switch trimmed[hashes] {
	case ' ', '\t':
		return hashes, true
	default:
		return 0, false
	}
}

func isHorizontalRule(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) < 3 {
		return false
	}

	for _, ch := range []byte{'-', '*', '_'} {
		same := true
		for i := 0; i < len(trimmed); i++ {
			if trimmed[i] != ch {
				same = false
				break
			}
		}
		if same {
			return true
		}
	}

	return false
}

func isTableRow(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}
	return strings.Contains(trimmed, "|")
}

func countMatches(pattern *regexp.Regexp, body string) int {
	return len(pattern.FindAllStringIndex(body, -1))
}

func countLinks(body string) int {
	count := 0
	for _, pattern := range []*regexp.Regexp{markdownLinkPattern, markdownLinkRefPattern} {
		for _, match := range pattern.FindAllStringIndex(body, -1) {
			if match[0] > 0 && body[match[0]-1] == '!' {
				continue
			}
			count++
		}
	}
	return count
}

func countInlineCodeSpans(body string) int {
	count := 0
	for i := 0; i < len(body); {
		if body[i] != '`' {
			i++
			continue
		}

		ticks := 1
		for i+ticks < len(body) && body[i+ticks] == '`' {
			ticks++
		}

		end := strings.Index(body[i+ticks:], strings.Repeat("`", ticks))
		if end == -1 {
			break
		}

		count++
		i += ticks + end + ticks
	}

	return count
}

var inlineCodePattern = regexp.MustCompile("`+[^`]*`+")

func countWords(text string) int {
	cleaned := stripForWordCount(text)
	if cleaned == "" {
		return 0
	}
	return len(strings.Fields(cleaned))
}

func stripForWordCount(text string) string {
	text = imagePattern.ReplaceAllStringFunc(text, altTextReplacer)
	text = imageRefPattern.ReplaceAllStringFunc(text, altTextReplacer)
	text = markdownLinkPattern.ReplaceAllStringFunc(text, altTextReplacer)
	text = markdownLinkRefPattern.ReplaceAllStringFunc(text, altTextReplacer)
	text = footnoteRefPattern.ReplaceAllString(text, " ")
	text = inlineCodePattern.ReplaceAllString(text, " ")

	var b strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			continue
		}
		if unicode.IsSpace(r) {
			b.WriteRune(' ')
		}
	}
	return strings.TrimSpace(b.String())
}

func altTextReplacer(match string) string {
	start := strings.Index(match, "[")
	end := strings.Index(match, "]")
	if start == -1 || end == -1 || end <= start+1 {
		return " "
	}
	return match[start+1:end] + " "
}

func countSentences(text string) int {
	cleaned := stripForSentenceCount(text)
	if cleaned == "" {
		return 0
	}

	count := 0
	for _, r := range cleaned {
		switch r {
		case '.', '!', '?':
			count++
		}
	}
	if count == 0 {
		return 1
	}
	return count
}

func stripForSentenceCount(text string) string {
	text = imagePattern.ReplaceAllStringFunc(text, altTextReplacer)
	text = imageRefPattern.ReplaceAllStringFunc(text, altTextReplacer)
	text = markdownLinkPattern.ReplaceAllStringFunc(text, altTextReplacer)
	text = markdownLinkRefPattern.ReplaceAllStringFunc(text, altTextReplacer)
	text = footnoteRefPattern.ReplaceAllString(text, " ")
	text = inlineCodePattern.ReplaceAllString(text, " ")

	var b strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '!' || r == '?' {
			b.WriteRune(r)
			continue
		}
		if unicode.IsSpace(r) {
			b.WriteRune(' ')
		}
	}
	return strings.TrimSpace(b.String())
}
