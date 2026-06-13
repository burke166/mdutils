package mdlint

const (
	SeverityError   = "error"
	SeverityWarning = "warning"

	RuleSingleH1               = "single-h1"
	RuleNoMissingH1            = "no-missing-h1"
	RuleNoSkippedHeadingLevels = "no-skipped-heading-levels"
	RuleNoDuplicateHeadings    = "no-duplicate-headings"
	RuleNoEmptyHeadings        = "no-empty-headings"
	RuleNoEmptySections        = "no-empty-sections"
	RuleNoTrailingWhitespace   = "no-trailing-whitespace"
	RuleMaxHeadingLength       = "max-heading-length"
	RuleMaxLineLength          = "max-line-length"
	RuleNoMultipleBlankLines   = "no-multiple-blank-lines"
	RuleRequireCodeFenceLang   = "require-code-fence-language"
	RuleNoEmptyLinks           = "no-empty-links"
	RuleRequireImageAltText    = "require-image-alt-text"
)

type Issue struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column,omitempty"`
	RuleID   string `json:"ruleId"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

func withSeverity(cfg Config, issue Issue) Issue {
	issue.Severity = cfg.SeverityFor(issue.RuleID)
	return issue
}

func exitCode(issues []Issue) int {
	hasError := false
	hasWarning := false

	for _, issue := range issues {
		switch issue.Severity {
		case SeverityError:
			hasError = true
		case SeverityWarning:
			hasWarning = true
		}
	}

	if hasError {
		return 2
	}
	if hasWarning {
		return 1
	}
	return 0
}
