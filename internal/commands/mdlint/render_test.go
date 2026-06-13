package mdlint

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderTextFormatsIssues(t *testing.T) {
	issues := []Issue{
		{
			File:     "README.md",
			Line:     1,
			RuleID:   RuleSingleH1,
			Severity: SeverityError,
			Message:  "document has multiple H1 headings",
		},
		{
			File:     "README.md",
			Line:     18,
			RuleID:   RuleNoSkippedHeadingLevels,
			Severity: SeverityWarning,
			Message:  "heading level skipped from H2 to H4",
		},
	}

	output := RenderText(issues, false)
	require.Contains(t, output, "README.md")
	require.Contains(t, output, "line 1: error single-h1: document has multiple H1 headings")
	require.Contains(t, output, "line 18: warning no-skipped-heading-levels:")
	require.Contains(t, output, "2 issues found")
}

func TestRenderTextQuietFiltersWarnings(t *testing.T) {
	issues := []Issue{
		{File: "a.md", Line: 1, RuleID: RuleSingleH1, Severity: SeverityError, Message: "error"},
		{File: "a.md", Line: 2, RuleID: RuleNoMissingH1, Severity: SeverityWarning, Message: "warn"},
	}

	output := RenderText(issues, true)
	require.Contains(t, output, "error")
	require.NotContains(t, output, "warn")
	require.Contains(t, output, "2 issues found")
}

func TestRenderJSON(t *testing.T) {
	issues := []Issue{
		{
			File:     "README.md",
			Line:     55,
			Column:   3,
			RuleID:   RuleNoDuplicateHeadings,
			Severity: SeverityWarning,
			Message:  `duplicate heading "Installation"`,
		},
	}

	output, err := RenderJSON(issues, false)
	require.NoError(t, err)
	require.Contains(t, output, `"file": "README.md"`)
	require.Contains(t, output, `"line": 55`)
	require.Contains(t, output, `"column": 3`)
	require.Contains(t, output, `"ruleId": "no-duplicate-headings"`)
	require.Contains(t, output, `"severity": "warning"`)
}

func TestRenderJSONQuiet(t *testing.T) {
	issues := []Issue{
		{File: "a.md", Line: 1, RuleID: RuleSingleH1, Severity: SeverityError, Message: "error"},
		{File: "a.md", Line: 2, RuleID: RuleNoMissingH1, Severity: SeverityWarning, Message: "warn"},
	}

	output, err := RenderJSON(issues, true)
	require.NoError(t, err)
	require.Contains(t, output, RuleSingleH1)
	require.NotContains(t, output, RuleNoMissingH1)
}

func TestFormatIssueLineWithColumn(t *testing.T) {
	line := formatIssueLine(Issue{
		Line:     10,
		Column:   4,
		Severity: SeverityWarning,
		RuleID:   RuleMaxLineLength,
		Message:  "too long",
	})
	require.Equal(t, "line 10:4: warning max-line-length: too long", line)
}

func TestRenderTextEmpty(t *testing.T) {
	require.Empty(t, strings.TrimSpace(RenderText(nil, false)))
}
