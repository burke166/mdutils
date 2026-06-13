package markdown

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitLines(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
	}{
		{
			name:    "empty content",
			content: "",
			want:    nil,
		},
		{
			name:    "single line without newline",
			content: "hello",
			want:    []string{"hello"},
		},
		{
			name:    "multiple lines",
			content: "a\nb\n",
			want:    []string{"a\n", "b\n", ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, SplitLines(tt.content))
		})
	}
}

func TestIsFenceLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		wantFence bool
		wantChar  byte
	}{
		{name: "backtick fence", line: "```go", wantFence: true, wantChar: '`'},
		{name: "tilde fence", line: "~~~", wantFence: true, wantChar: '~'},
		{name: "indented fence", line: "  ```", wantFence: true, wantChar: '`'},
		{name: "plain text", line: "not a fence", wantFence: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFence, gotChar := IsFenceLine(tt.line)
			require.Equal(t, tt.wantFence, gotFence)
			require.Equal(t, tt.wantChar, gotChar)
		})
	}
}

func TestIsBlankLine(t *testing.T) {
	tests := []struct {
		name string
		line string
		want bool
	}{
		{name: "empty", line: "", want: true},
		{name: "spaces", line: "   \n", want: true},
		{name: "text", line: "text\n", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, IsBlankLine(tt.line))
		})
	}
}

func TestTrailingNewline(t *testing.T) {
	tests := []struct {
		name string
		line string
		want string
	}{
		{name: "lf", line: "a\n", want: "\n"},
		{name: "crlf", line: "a\r\n", want: "\r\n"},
		{name: "none", line: "a", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, TrailingNewline(tt.line))
		})
	}
}
