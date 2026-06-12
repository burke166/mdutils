package output

import (
	"strings"
	"testing"
)

func TestRenderNumbered(t *testing.T) {
	actual := strings.TrimSpace(RenderNumbered(testHeadings))

	expected := strings.TrimSpace(`
1. Adventure Wargame
  1.1. Character Creation
    1.1.1. Attributes
  1.2. Equipment
`)

	if actual != expected {
		t.Errorf("unexpected numbered output\n\nexpected:\n%s\n\nactual:\n%s", expected, actual)
	}
}
