package subtitle

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

// Check error and panic
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// The Split function that detects a *subtitle block* in a SRT file
func SplitSubtitles(data []byte, atEOF bool) (advance int, token []byte, err error) {

	// Return nothing if at end of file and no data passed
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// If at end of file with data return the data
	if atEOF {
		return len(data), data, nil
	}

	// This is the separator string
	sep := regexp.MustCompile(`(\s*\r?\n){2,}`)

	// Find the index of the sep regexp in data
	loc := sep.FindIndex(data)
	if loc == nil {
		// Not found
		// return len(data), data, nil
		return 0, nil, nil
	} else {
		return loc[1], data[:loc[0]], nil
	}
}

// prepare a string, clean up, etc.
func PrepareString(data string) string {
	// convert []byte to string
	text := string(data)

	// Note that \s == [ \t\f\n\r\v]

	// Substitute double quotes to single quotes ?
	text = regexp.MustCompile(`"`).ReplaceAllString(text, "'")
	// Substitute multiple spaces to single space
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	// Change space+punctuation to puntuation alone
	text = regexp.MustCompile(`\s([,:;!\\?\.\)\]])`).ReplaceAllString(text, "$1")
	// Delete spaces at the beginning and end of text
	text = strings.TrimSpace(text)

	return text
}

func JoinAllStrings(data ...string) string {
	text := new(strings.Builder)

	// Create a string adding all the lines
	for _, theLine := range data {
		text.WriteString(theLine)
		text.WriteString(" ")
	}
	return strings.TrimSpace(text.String())
}

// Find closest extreme to the center
// returns true if the first is closest or false if the second
func ClosestFloat(center, first, second float64) bool {
	if math.Abs(center-first) <= math.Abs(center-second) {
		return true
	} else {
		return false
	}
}

// Returns the closest string and lenght
func ClosestNotEmptyString(center float64, str1, str2 string) (string, float64) {
	len1 := float64(len([]rune(str1)))
	len2 := float64(len([]rune(str2)))
	if ClosestFloat(center, len1, len2) && str1 != "" {
		return str1, len1 - center
	}
	return str2, len2 - center
}

func ConcatWithSpace(str1, str2 string) string {
	if str1 == "[]" || str1 == "" {
		return str2
	}
	if str2 == "[]" || str2 == "" {
		return str1
	}
	return str1 + " " + str2
}

func PrintStringMaxWidth(s string, width int) string {

	lenString := len([]rune(s))

	// Print the text
	if lenString > width {
		return fmt.Sprintf("%-*.*s...", width-3, width-3, s)
	} else {
		return fmt.Sprintf("%-*.*s", width, width, s)
	}
}
