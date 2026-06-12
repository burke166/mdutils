package output

import (
	"encoding/json"
	"testing"

	"github.com/computercodeblue/mdutils/internal/markdown"
	"github.com/stretchr/testify/require"
)

func TestRenderJson(t *testing.T) {
	actual, err := RenderJson(testHeadings)
	require.NoError(t, err)

	var got []markdown.Heading
	err = json.Unmarshal([]byte(actual), &got)
	require.NoError(t, err)

	require.Equal(t, testHeadings, got)
}
