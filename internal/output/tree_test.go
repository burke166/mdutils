package output

import (
	"strings"
	"testing"
)

func TestRenderTree(t *testing.T) {
	actual := strings.TrimSpace(RenderTree(testHeadings))

	expected := strings.TrimSpace(`
Adventure Wargame
├── Character Creation
│   └── Attributes
└── Equipment
`)

	if actual != expected {
		t.Errorf("unexpected tree output\n\nexpected:\n%s\n\nactual:\n%s", expected, actual)
	}
}
