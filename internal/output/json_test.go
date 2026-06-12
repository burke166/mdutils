package output

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

func TestRenderJson(t *testing.T) {
	actual, err := RenderJson(testHeadings)
	if err != nil {
		t.Fatalf("unable to render JSON: %v", err)
	}

	var got []markdown.Heading
	if err := json.Unmarshal([]byte(actual), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if !reflect.DeepEqual(got, testHeadings) {
		t.Errorf("unexpected JSON output\n\ngot: %#v\n\nwant: %#v", got, testHeadings)
	}
}
