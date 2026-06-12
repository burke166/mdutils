package output

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderTree(t *testing.T) {
	actual := strings.TrimSpace(RenderTree(testHeadings))

	expected := strings.TrimSpace(`
Adventure Wargame
├── Character Creation
│   └── Attributes
└── Equipment
`)

	require.Equal(t, expected, actual)
}
