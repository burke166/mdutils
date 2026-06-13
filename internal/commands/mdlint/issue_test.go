package mdlint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExitCodeNoIssues(t *testing.T) {
	require.Equal(t, 0, exitCode(nil))
}

func TestExitCodeWarningsOnly(t *testing.T) {
	issues := []Issue{{Severity: SeverityWarning}}
	require.Equal(t, 1, exitCode(issues))
}

func TestExitCodeErrors(t *testing.T) {
	issues := []Issue{
		{Severity: SeverityWarning},
		{Severity: SeverityError},
	}
	require.Equal(t, 2, exitCode(issues))
}

func TestWithSeverityUsesConfig(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Severity[RuleNoMissingH1] = SeverityError

	issue := withSeverity(cfg, Issue{RuleID: RuleNoMissingH1})
	require.Equal(t, SeverityError, issue.Severity)
}
