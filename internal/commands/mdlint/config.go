package mdlint

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const configFileName = ".mdlintrc.yaml"

type Config struct {
	Rules    RuleConfig        `yaml:"rules"`
	Severity map[string]string `yaml:"severity"`
	Exclude  []string          `yaml:"exclude"`
}

type RuleConfig struct {
	SingleH1                 bool `yaml:"single-h1"`
	NoMissingH1              bool `yaml:"no-missing-h1"`
	NoSkippedHeadingLevels   bool `yaml:"no-skipped-heading-levels"`
	NoDuplicateHeadings      bool `yaml:"no-duplicate-headings"`
	NoEmptyHeadings          bool `yaml:"no-empty-headings"`
	NoEmptySections          bool `yaml:"no-empty-sections"`
	NoTrailingWhitespace     bool `yaml:"no-trailing-whitespace"`
	MaxHeadingLength         int  `yaml:"max-heading-length"`
	MaxLineLength            int  `yaml:"max-line-length"`
	NoMultipleBlankLines     bool `yaml:"no-multiple-blank-lines"`
	RequireCodeFenceLanguage bool `yaml:"require-code-fence-language"`
	NoEmptyLinks             bool `yaml:"no-empty-links"`
	RequireImageAltText      bool `yaml:"require-image-alt-text"`
}

func DefaultConfig() Config {
	return Config{
		Rules: RuleConfig{
			SingleH1:                 true,
			NoMissingH1:              true,
			NoSkippedHeadingLevels:   true,
			NoDuplicateHeadings:      true,
			NoEmptyHeadings:          true,
			NoEmptySections:          true,
			NoTrailingWhitespace:     true,
			MaxHeadingLength:         80,
			MaxLineLength:            120,
			NoMultipleBlankLines:     true,
			RequireCodeFenceLanguage: false,
			NoEmptyLinks:             true,
			RequireImageAltText:      false,
		},
		Severity: map[string]string{
			RuleSingleH1:               SeverityError,
			RuleNoMissingH1:            SeverityWarning,
			RuleNoSkippedHeadingLevels: SeverityWarning,
			RuleNoDuplicateHeadings:    SeverityWarning,
			RuleNoEmptyHeadings:        SeverityError,
			RuleNoEmptySections:        SeverityWarning,
			RuleNoTrailingWhitespace:   SeverityWarning,
			RuleMaxHeadingLength:       SeverityWarning,
			RuleMaxLineLength:          SeverityWarning,
			RuleNoMultipleBlankLines:   SeverityWarning,
			RuleRequireCodeFenceLang:   SeverityWarning,
			RuleNoEmptyLinks:           SeverityError,
			RuleRequireImageAltText:    SeverityWarning,
		},
		Exclude: nil,
	}
}

func LoadConfig(explicitPath string) (Config, error) {
	if explicitPath != "" {
		return loadConfigFile(explicitPath)
	}

	for _, path := range configSearchPaths() {
		if path == "" {
			continue
		}
		if _, err := os.Stat(path); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return Config{}, err
		}
		return loadConfigFile(path)
	}

	return DefaultConfig(), nil
}

func loadConfigFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config %s: %w", path, err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config %s: %w", path, err)
	}

	return cfg, nil
}

func configSearchPaths() []string {
	var paths []string

	if cwd, err := os.Getwd(); err == nil {
		paths = append(paths, filepath.Join(cwd, configFileName))
	}

	if home, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(home, configFileName))
	}

	if exeDir, err := executableDir(); err == nil {
		paths = append(paths, filepath.Join(exeDir, configFileName))
	}

	return paths
}

func executableDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

func (c Config) SeverityFor(ruleID string) string {
	if severity, ok := c.Severity[ruleID]; ok && severity != "" {
		return severity
	}
	return SeverityWarning
}
