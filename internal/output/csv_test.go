package output

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderCsv(t *testing.T) {
	actual, err := RenderCsv(testHeadings)
	require.NoError(t, err)
	actual = strings.TrimSpace(actual)

	expected := strings.TrimSpace(`
level,text
1,Adventure Wargame
2,Character Creation
3,Attributes
2,Equipment
`)

	require.Equal(t, expected, actual)
}
