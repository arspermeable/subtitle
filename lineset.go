package subtitle

import (
	"strings"
)

// A LineSet is a set of lines within the list of subtitle text lines
// Each LineSet is process as a block to match original and translation linebreaks
type LineSet struct {
	initLine       int
	lastLine       int
	translatedText string
}

func JoinAllStringsFromLineSet(data []LineSet) string {
	text := new(strings.Builder)
	var init bool = true

	// Create a string adding all the lines
	for _, theLineSet := range data {
		var err error
		if !init {
			_, err = text.WriteString(" ")
			check(err)
		}
		_, err = text.WriteString(theLineSet.translatedText)
		check(err)
		init = false
	}
	return text.String()
}
