package output

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderedHeadings(t *testing.T) {
	actual := strings.TrimSpace(RenderMarkdownHeadings(testHeadings))

	expected := strings.TrimSpace(`
# Adventure Wargame
## Character Creation
### Attributes
## Equipment
`)

	require.Equal(t, expected, actual)
}
