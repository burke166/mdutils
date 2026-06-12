package output

import (
	"strings"
	"testing"
)

func TestRenderCsv(t *testing.T) {
	actual, err := RenderCsv(testHeadings)
	if err != nil {
		t.Fatalf("unable to render CSV: %v", err)
	}
	actual = strings.TrimSpace(actual)

	expected := strings.TrimSpace(`
level,text
1,Adventure Wargame
2,Character Creation
3,Attributes
2,Equipment
`)

	if actual != expected {
		t.Errorf("unexpected CSV output\n\nexpected:\n%s\n\nactual:\n%s", expected, actual)
	}
}
