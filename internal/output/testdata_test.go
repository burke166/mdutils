package output

import "github.com/computercodeblue/mdutils/internal/markdown"

var testHeadings = []markdown.Heading{
	{Level: 1, Text: "Adventure Wargame"},
	{Level: 2, Text: "Character Creation"},
	{Level: 3, Text: "Attributes"},
	{Level: 2, Text: "Equipment"},
}

var emptyHeadings = []markdown.Heading{}

var singleHeading = []markdown.Heading{
	{Level: 1, Text: "Single Heading"},
}
