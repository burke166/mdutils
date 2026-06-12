package mdcheck

import (
	"testing"

	"github.com/computercodeblue/mdutils/internal/markdown"
	"github.com/stretchr/testify/require"
)

func defaultValidateOptions() ValidateOptions {
	return ValidateOptions{
		CheckDuplicates: true,
	}
}

func extractHeadings(t *testing.T, source string) []markdown.Heading {
	t.Helper()

	headings, err := markdown.ExtractHeadings([]byte(source))
	require.NoError(t, err)

	return headings
}

func TestValidateMissingH1(t *testing.T) {
	headings := extractHeadings(t, "## Section\n")
	diags := Validate(headings, defaultValidateOptions())

	require.Len(t, diags, 1)
	require.Equal(t, RuleMissingH1, diags[0].Rule)
	require.Equal(t, "document has no H1 heading", diags[0].Message)
}

func TestValidateMissingH1NotReportedWhenPresent(t *testing.T) {
	headings := extractHeadings(t, "# Title\n\n## Section\n")
	diags := Validate(headings, defaultValidateOptions())

	for _, d := range diags {
		require.NotEqual(t, RuleMissingH1, d.Rule)
	}
}

func TestValidateMultipleH1(t *testing.T) {
	headings := extractHeadings(t, "# First\n\n# Second\n\n# Third\n")
	diags := Validate(headings, defaultValidateOptions())

	var multiple []Diagnostic
	for _, d := range diags {
		if d.Rule == RuleMultipleH1 {
			multiple = append(multiple, d)
		}
	}

	require.Len(t, multiple, 2)
	require.Equal(t, "Second", multiple[0].Text)
	require.Equal(t, "Third", multiple[1].Text)
}

func TestValidateAllowMultipleH1(t *testing.T) {
	headings := extractHeadings(t, "# First\n\n# Second\n")
	diags := Validate(headings, ValidateOptions{
		CheckDuplicates: true,
		AllowMultipleH1: true,
	})

	for _, d := range diags {
		require.NotEqual(t, RuleMultipleH1, d.Rule)
	}
}

func TestValidateSkippedLevel(t *testing.T) {
	headings := extractHeadings(t, "# Title\n\n## Section\n\n#### Deep\n")
	diags := Validate(headings, defaultValidateOptions())

	require.Len(t, diags, 1)
	require.Equal(t, RuleSkippedLevel, diags[0].Rule)
	require.Equal(t, 5, diags[0].Line)
	require.Equal(t, 4, diags[0].Level)
	require.Equal(t, "Deep", diags[0].Text)
	require.Contains(t, diags[0].Message, "H2 followed by H4")
}

func TestValidateSkippedLevelNotReportedForValidNesting(t *testing.T) {
	headings := extractHeadings(t, "# Title\n\n## Section\n\n### Subsection\n\n## Another\n")
	diags := Validate(headings, defaultValidateOptions())

	for _, d := range diags {
		require.NotEqual(t, RuleSkippedLevel, d.Rule)
	}
}

func TestValidateDuplicateHeading(t *testing.T) {
	headings := extractHeadings(t, "# Title\n\n## Section\n\n## Section\n")
	diags := Validate(headings, defaultValidateOptions())

	require.Len(t, diags, 1)
	require.Equal(t, RuleDuplicateHeading, diags[0].Rule)
	require.Equal(t, "Section", diags[0].Text)
}

func TestValidateNoDuplicates(t *testing.T) {
	headings := extractHeadings(t, "# Title\n\n## Section\n\n## Section\n")
	diags := Validate(headings, ValidateOptions{
		CheckDuplicates: false,
	})

	for _, d := range diags {
		require.NotEqual(t, RuleDuplicateHeading, d.Rule)
	}
}

func TestValidateEmptyHeading(t *testing.T) {
	headings := extractHeadings(t, "# Title\n\n## \n")
	diags := Validate(headings, defaultValidateOptions())

	require.Len(t, diags, 1)
	require.Equal(t, RuleEmptyHeading, diags[0].Rule)
	require.Equal(t, 3, diags[0].Line)
	require.Equal(t, 2, diags[0].Level)
}

func TestValidateMaxLevel(t *testing.T) {
	headings := extractHeadings(t, "# Title\n\n## Section\n\n### Subsection\n\n#### Deep\n")
	diags := Validate(headings, ValidateOptions{
		CheckDuplicates: true,
		MaxLevel:        3,
	})

	require.Len(t, diags, 1)
	require.Equal(t, RuleMaxLevel, diags[0].Rule)
	require.Equal(t, 7, diags[0].Line)
	require.Equal(t, 4, diags[0].Level)
	require.Equal(t, "Deep", diags[0].Text)
}

func TestValidateMaxLevelDisabledWhenZero(t *testing.T) {
	headings := extractHeadings(t, "# Title\n\n###### Deep\n")
	diags := Validate(headings, ValidateOptions{
		CheckDuplicates: true,
		MaxLevel:        0,
	})

	for _, d := range diags {
		require.NotEqual(t, RuleMaxLevel, d.Rule)
	}
}

func TestRenderHuman(t *testing.T) {
	output := RenderHuman([]Diagnostic{{
		Rule:    RuleSkippedLevel,
		Message: "skipped heading level (H2 followed by H4)",
		Line:    12,
		Level:   4,
		Text:    "Details",
	}}, false)

	require.Equal(t, `line 12: [skipped-level] H4 "Details": skipped heading level (H2 followed by H4)`, output)
}

func TestRenderHumanMissingH1(t *testing.T) {
	output := RenderHuman([]Diagnostic{{
		Rule:    RuleMissingH1,
		Message: "document has no H1 heading",
	}}, false)

	require.Equal(t, "document: [missing-h1] document has no H1 heading", output)
}

func TestRenderJSON(t *testing.T) {
	output, err := RenderJSON([]Diagnostic{{
		Rule:    RuleEmptyHeading,
		Message: "empty heading",
		Line:    3,
		Level:   2,
	}})
	require.NoError(t, err)

	require.JSONEq(t, `[
  {
    "rule": "empty-heading",
    "message": "empty heading",
    "line": 3,
    "level": 2,
    "text": ""
  }
]`, output)
}

func TestRenderJSONEmpty(t *testing.T) {
	output, err := RenderJSON(nil)
	require.NoError(t, err)
	require.Equal(t, "[]\n", output)
}

func TestOrderDiagnosticsByLine(t *testing.T) {
	diags := []Diagnostic{
		{Rule: RuleMultipleH1, Line: 5, Level: 1, Text: "Second", Message: "document has more than one H1 heading"},
		{Rule: RuleSkippedLevel, Line: 3, Level: 4, Text: "Deep", Message: "skipped heading level (H1 followed by H4)"},
	}

	ordered := orderDiagnostics(diags, false)
	require.Equal(t, RuleSkippedLevel, ordered[0].Rule)
	require.Equal(t, RuleMultipleH1, ordered[1].Rule)
}

func TestOrderDiagnosticsByRule(t *testing.T) {
	diags := []Diagnostic{
		{Rule: RuleMultipleH1, Line: 5, Level: 1, Text: "Second", Message: "document has more than one H1 heading"},
		{Rule: RuleSkippedLevel, Line: 3, Level: 4, Text: "Deep", Message: "skipped heading level (H1 followed by H4)"},
	}

	ordered := orderDiagnostics(diags, true)
	require.Equal(t, RuleMultipleH1, ordered[0].Rule)
	require.Equal(t, RuleSkippedLevel, ordered[1].Rule)
}
