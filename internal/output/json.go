package output

import (
	"encoding/json"

	"github.com/computercodeblue/mdutils/internal/markdown"
)

func RenderJson(headings []markdown.Heading) (string, error) {
	data, err := json.MarshalIndent(headings, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data) + "\n", nil
}
