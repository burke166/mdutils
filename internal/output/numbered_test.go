package output

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderNumbered(t *testing.T) {
	actual := strings.TrimSpace(RenderNumbered(testHeadings))

	expected := strings.TrimSpace(`
1. Adventure Wargame
  1.1. Character Creation
    1.1.1. Attributes
  1.2. Equipment
`)

	require.Equal(t, expected, actual)
}
