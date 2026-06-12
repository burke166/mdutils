package output

import (
	"strings"
	"testing"
)

func TestRenderedHeadings(t *testing.T) {
	actual := strings.TrimSpace(RenderMarkdownHeadings(testHeadings))

	expected := strings.TrimSpace(`
# Adventure Wargame
## Character Creation
### Attributes
## Equipment
`)

	if actual != expected {
		t.Errorf("unexpected headings output\n\nexpected:\n%s\n\nactual:\n%s", expected, actual)
	}
}
