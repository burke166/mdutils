package mdlint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	require.True(t, cfg.Rules.SingleH1)
	require.True(t, cfg.Rules.NoMissingH1)
	require.Equal(t, 80, cfg.Rules.MaxHeadingLength)
	require.Equal(t, 120, cfg.Rules.MaxLineLength)
	require.False(t, cfg.Rules.RequireCodeFenceLanguage)
	require.Equal(t, SeverityError, cfg.Severity[RuleSingleH1])
	require.Equal(t, SeverityWarning, cfg.Severity[RuleNoMissingH1])
}

func TestLoadConfigExplicitPath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "custom.yaml")

	content := `
rules:
  single-h1: false
  max-line-length: 90
severity:
  single-h1: warning
exclude:
  - "notes.md"
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	cfg, err := LoadConfig(path)
	require.NoError(t, err)
	require.False(t, cfg.Rules.SingleH1)
	require.Equal(t, 90, cfg.Rules.MaxLineLength)
	require.Equal(t, SeverityWarning, cfg.Severity[RuleSingleH1])
	require.Equal(t, []string{"notes.md"}, cfg.Exclude)
}

func TestLoadConfigMissingExplicitPath(t *testing.T) {
	_, err := LoadConfig(filepath.Join(t.TempDir(), "missing.yaml"))
	require.Error(t, err)
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, configFileName)
	require.NoError(t, os.WriteFile(path, []byte("rules: ["), 0644))

	_, err := LoadConfig(path)
	require.Error(t, err)
}

func TestLoadConfigDiscoveryOrderUsesCWD(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)

	cwd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})

	path := filepath.Join(dir, configFileName)
	content := `
rules:
  single-h1: false
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	cfg, err := LoadConfig("")
	require.NoError(t, err)
	require.False(t, cfg.Rules.SingleH1)
}

func TestLoadConfigDiscoveryOrderUsesHome(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)

	work := t.TempDir()
	cwd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(work))
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})

	content := `
rules:
  max-line-length: 42
`
	require.NoError(t, os.WriteFile(filepath.Join(home, configFileName), []byte(content), 0644))

	cfg, err := LoadConfig("")
	require.NoError(t, err)
	require.Equal(t, 42, cfg.Rules.MaxLineLength)
}

func TestLoadConfigFallsBackToDefaults(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)

	cwd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})

	cfg, err := LoadConfig("")
	require.NoError(t, err)
	require.Equal(t, DefaultConfig(), cfg)
}

func TestLoadConfigPartialMergeKeepsDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, configFileName)
	content := `
rules:
  require-code-fence-language: true
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	cfg, err := loadConfigFile(path)
	require.NoError(t, err)
	require.True(t, cfg.Rules.RequireCodeFenceLanguage)
	require.True(t, cfg.Rules.SingleH1)
	require.Equal(t, DefaultConfig().Severity, cfg.Severity)
}

func TestSeverityForUnknownRule(t *testing.T) {
	cfg := DefaultConfig()
	require.Equal(t, SeverityWarning, cfg.SeverityFor("unknown-rule"))
}
