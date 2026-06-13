package mdlint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func testConfig() Config {
	return DefaultConfig()
}

func lintString(t *testing.T, content string, cfg Config) []Issue {
	t.Helper()
	issues, err := LintContent("test.md", []byte(content), cfg)
	require.NoError(t, err)
	return issues
}

func hasRule(issues []Issue, ruleID string) bool {
	for _, issue := range issues {
		if issue.RuleID == ruleID {
			return true
		}
	}
	return false
}

func TestLintValidDocument(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "simple.md")
	issues, err := LintFile(path, testConfig())
	require.NoError(t, err)
	require.Empty(t, issues)
}

func TestRuleNoMissingH1(t *testing.T) {
	issues := lintString(t, "## Section\n", testConfig())
	require.True(t, hasRule(issues, RuleNoMissingH1))
}

func TestRuleSingleH1(t *testing.T) {
	issues := lintString(t, "# First\n\n# Second\n", testConfig())
	require.True(t, hasRule(issues, RuleSingleH1))
}

func TestRuleSingleH1Disabled(t *testing.T) {
	cfg := testConfig()
	cfg.Rules.SingleH1 = false

	issues := lintString(t, "# First\n\n# Second\n", cfg)
	require.False(t, hasRule(issues, RuleSingleH1))
}

func TestRuleNoSkippedHeadingLevels(t *testing.T) {
	issues := lintString(t, "# Title\n\n## Section\n\n#### Deep\n", testConfig())
	require.True(t, hasRule(issues, RuleNoSkippedHeadingLevels))
}

func TestRuleNoDuplicateHeadings(t *testing.T) {
	issues := lintString(t, "# Title\n\n## Installation\n\n## Installation\n", testConfig())
	require.True(t, hasRule(issues, RuleNoDuplicateHeadings))
}

func TestRuleNoEmptyHeadings(t *testing.T) {
	issues := lintString(t, "# Title\n\n## \n", testConfig())
	require.True(t, hasRule(issues, RuleNoEmptyHeadings))
}

func TestRuleNoEmptySections(t *testing.T) {
	issues := lintString(t, "# Title\n\n## Empty\n\n## Next\n\nContent\n", testConfig())
	require.True(t, hasRule(issues, RuleNoEmptySections))
}

func TestRuleMaxHeadingLength(t *testing.T) {
	long := strings.Repeat("a", 81)
	issues := lintString(t, "# "+long+"\n", testConfig())
	require.True(t, hasRule(issues, RuleMaxHeadingLength))
}

func TestRuleNoTrailingWhitespace(t *testing.T) {
	issues := lintString(t, "# Title  \n", testConfig())
	require.True(t, hasRule(issues, RuleNoTrailingWhitespace))
}

func TestRuleMaxLineLength(t *testing.T) {
	long := strings.Repeat("a", 121)
	issues := lintString(t, long+"\n", testConfig())
	require.True(t, hasRule(issues, RuleMaxLineLength))
}

func TestRuleNoMultipleBlankLines(t *testing.T) {
	issues := lintString(t, "# Title\n\n\nBody\n", testConfig())
	require.True(t, hasRule(issues, RuleNoMultipleBlankLines))
}

func TestRuleRequireCodeFenceLanguage(t *testing.T) {
	cfg := testConfig()
	cfg.Rules.RequireCodeFenceLanguage = true

	issues := lintString(t, "# Title\n\n```\ncode\n```\n", cfg)
	require.True(t, hasRule(issues, RuleRequireCodeFenceLang))
}

func TestRuleRequireCodeFenceLanguageIgnoredWithLang(t *testing.T) {
	cfg := testConfig()
	cfg.Rules.RequireCodeFenceLanguage = true

	issues := lintString(t, "# Title\n\n```go\ncode\n```\n", cfg)
	require.False(t, hasRule(issues, RuleRequireCodeFenceLang))
}

func TestRuleNoEmptyLinks(t *testing.T) {
	issues := lintString(t, "# Title\n\n[text]()\n", testConfig())
	require.True(t, hasRule(issues, RuleNoEmptyLinks))

	issues = lintString(t, "# Title\n\n[]()\n", testConfig())
	require.True(t, hasRule(issues, RuleNoEmptyLinks))
}

func TestRuleRequireImageAltText(t *testing.T) {
	cfg := testConfig()
	cfg.Rules.RequireImageAltText = true

	issues := lintString(t, "# Title\n\n![](image.png)\n", cfg)
	require.True(t, hasRule(issues, RuleRequireImageAltText))
}

func TestHeadingsInsideCodeBlocksIgnored(t *testing.T) {
	path := filepath.Join("..", "..", "..", "testdata", "codeblocks.md")
	issues, err := LintFile(path, testConfig())
	require.NoError(t, err)
	require.False(t, hasRule(issues, RuleSingleH1))
	require.False(t, hasRule(issues, RuleNoDuplicateHeadings))
}

func TestLintFileReadsFromDisk(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "doc.md")
	require.NoError(t, os.WriteFile(path, []byte("## Missing H1\n"), 0644))

	issues, err := LintFile(path, testConfig())
	require.NoError(t, err)
	require.True(t, hasRule(issues, RuleNoMissingH1))
}
