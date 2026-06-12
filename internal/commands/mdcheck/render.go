package mdcheck

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

func RenderHuman(diagnostics []Diagnostic, groupByRule bool) string {
	if len(diagnostics) == 0 {
		return ""
	}

	diagnostics = orderDiagnostics(diagnostics, groupByRule)

	var b strings.Builder
	for i, d := range diagnostics {
		if i > 0 {
			b.WriteByte('\n')
		}

		b.WriteString(formatDiagnostic(d))
	}

	return b.String()
}

func orderDiagnostics(diagnostics []Diagnostic, groupByRule bool) []Diagnostic {
	ordered := append([]Diagnostic(nil), diagnostics...)
	if groupByRule {
		return ordered
	}

	slices.SortFunc(ordered, func(a, b Diagnostic) int {
		if a.Line != b.Line {
			return a.Line - b.Line
		}

		return strings.Compare(a.Rule, b.Rule)
	})

	return ordered
}

func RenderJSON(diagnostics []Diagnostic) (string, error) {
	if diagnostics == nil {
		diagnostics = []Diagnostic{}
	}

	data, err := json.MarshalIndent(diagnostics, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data) + "\n", nil
}

func formatDiagnostic(d Diagnostic) string {
	location := "document"
	if d.Line > 0 {
		location = fmt.Sprintf("line %d", d.Line)
	}

	heading := headingLabel(d.Level, d.Text)
	if heading != "" {
		return fmt.Sprintf("%s: [%s] %s: %s", location, d.Rule, heading, d.Message)
	}

	return fmt.Sprintf("%s: [%s] %s", location, d.Rule, d.Message)
}

func headingLabel(level int, text string) string {
	if level <= 0 {
		return ""
	}

	if text == "" {
		return fmt.Sprintf("H%d", level)
	}

	return fmt.Sprintf("H%d %q", level, text)
}
