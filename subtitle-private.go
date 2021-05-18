package subtitle

import (
	"fmt"
	"regexp"
	"strings"
)

// Add a subtitle block to the SubtitleSRT
func (this *SubtitleSRT) appendSubtitle(data string) {

	// Split the data by any newline
	lines := regexp.MustCompile(`\r?\n`).Split(data, -1)

	// Extract the lines
	nLine := 2
	for OK := true; OK; OK = nLine < len(lines) {
		theLine := "[]"
		// If there is a line, get it
		if len(lines) > nLine {
			theLine = PrepareString(lines[nLine])
		}
		// If it is not an empty line, or it is the first one, add it
		if theLine != "" || nLine == 2 {
			this.originalLine = append(this.originalLine, theLine)
			nLine++
		} else {
			break
		}
	}
	// Populate and append the subtitle data
	this.subtitleBlock = append(this.subtitleBlock, SubtitleBlock{lines[0], lines[1], nLine - 2})
}

// Split what is in translatedText, into line sets detected by
// comparing original lines with translatedText
func (this *SubtitleSRT) splitTranslatedTextIntoLineSets() {

	// Kind: isExact or not - initialy not
	isExact := false
	// Translated text to be splitted
	data := this.translatedText
	// Create the first newLineSet
	newLineSet := LineSet{0, 0}
	newTranslatedSet := ""
	// String for the searchRegexp
	searchRegexp := ""
	// Boolean to check if it isMiniLine
	isMiniLine := false

	for i, theLine := range this.originalLine {

		searchRegexp = ""
		isMiniLine = false

		// Set the searchRegExp
		if theLine == "" {
			// This line is empty, search for []
			searchRegexp = `\[\]`
		} else if matched, _ := regexp.MatchString(`\[.+\]`, theLine); matched {
			// This line is a word+ between brackets
			searchRegexp = `\[([^]]+)\]`
		} else {
			// This line is a text line
			searchRegexp = regexp.QuoteMeta(theLine)
			// Verify if it isMiniLine
			isMiniLine = len([]rune(theLine)) < minmatch
		}

		// Can the searchRegexp be found in the translation?
		loc := regexp.MustCompile(searchRegexp).FindStringIndex(data)

		// !found has three cases:
		//    - OR loc==nil (really not found)
		//    - OR found in the middle of the text, but it is a miniLine
		//    - OR regular line found in the middle of the text while processing an isExact set
		//			this is a weird case, because... where do we put the text in the middle?
		found := !(loc == nil || (loc[0] != 0 && isMiniLine) || (loc[0] != 0 && !isMiniLine && isExact))
		if !found {
			// The line was not found
			if isExact {
				// If current line set isExact, append the new line set
				this.lineSet = append(this.lineSet, newLineSet)
				this.translatedSet = append(this.translatedSet, newTranslatedSet)
				// and open a new lineset that !isExact
				newLineSet = LineSet{i, i}
				newTranslatedSet = ""
				isExact = false
			} else {
				// If newLineSet set !isExact, add this line to the newLineSet
				newLineSet.LastLine = i
			}
		} else {
			// The line was found
			if loc[0] == 0 {
				// The searchRegexp is at the beginning of the data
				// newLineSet only can be isExact, or it is the initial one
				// Add the line to the newLineSet and the text to the translated line
				newLineSet.LastLine = i
				newTranslatedSet = ConcatWithSpace(newTranslatedSet, data[loc[0]:loc[1]])
				// This line set continues isExact (just in case we are in the initial line)
				isExact = true
				// Retire the found string
				data = data[loc[1]+1:]
			} else {
				// The SearchRegexp is in the middle of the data.
				// It cannot be isMiniLine because that is in !found case
				// And it cannot be !isMiniLine and isExact for the same reason
				// So, a new isExact line set is found while treating a !isExact newLineSet
				// Add the text and close previous line set
				newTranslatedSet = data[:loc[0]-1]
				this.lineSet = append(this.lineSet, newLineSet)
				this.translatedSet = append(this.translatedSet, newTranslatedSet)
				// Open a newLineSet that isExact
				newLineSet = LineSet{i, i}
				newTranslatedSet = ConcatWithSpace("", data[loc[0]:loc[1]])
				isExact = true
				data = data[loc[1]+1:]
			}
		}
	}

	// If there is something left in data, it is the last block
	if data != "" {
		newTranslatedSet = strings.TrimSpace(data)
		this.lineSet = append(this.lineSet, newLineSet)
		this.translatedSet = append(this.translatedSet, newTranslatedSet)
	}
}

// Split the text assigned to a LineSet into lines
func (this *SubtitleSRT) splitTranslatedLineSetIntoLines(theLineSet int) {

	// Calculate the ratio translated:original for this line set
	ratio := this.CalculateRatioOfLineSet(theLineSet)

	// Initialize the values
	data := this.translatedSet[theLineSet]
	excess := 0.0

	// Iterate over the lines of the lineSet theLineSet
	for i := this.lineSet[theLineSet].InitLine; i <= this.lineSet[theLineSet].LastLine; i++ {
		// Set target size for the translated line
		// Note that target is float64!

		lenOrig := len([]rune(this.originalLine[i]))
		target := ratio*float64(lenOrig) - excess
		newLine := ""

		if this.CountOriginalCharsInLine(i) == 0 {
			// If the original line is empty, output empty data and keep excess
			newLine = ""
		} else if i == this.lineSet[theLineSet].LastLine {
			// If this is the last line, output the rest of the data
			newLine = data
			excess = float64(len([]rune(newLine))) - target + 1.0
		} else if data == "" {
			// if the remaining string is "", return an empty string
			newLine = ""
			excess = -target
			// and raise an error here!!
		} else {
			var chars int
			var subStrMax, subStrMin, re string
			// Find the next "split point" after "target" characters, defined as
			// <chars>{target-1} + <no-sep>{+} + <sep>{+}
			if target > 1 {
				chars = int(target + 0.5)
			}
			// Get the substring until the next separator after {chars}
			// (****)!! Note: It must be {chars} unicode points not ascii (.)
			re = fmt.Sprintf(`^([\P{M}\p{M}]{%d}[^%s]*[%s]*)`, chars, sep, sep)
			subStrMax = PrepareString(regexp.MustCompile(re).FindString(data))
			if subStrMax == "" {
				subStrMax = data
			}
			// Get the substring until the prev separator before .{chars}
			// re = fmt.Sprintf("([%s]*[^%s]*[%s]*)$", sep, sep, sep)
			re = fmt.Sprintf("(\\p{Z}*[^%s]*[%s]*)$", sep, sep)
			subStrMin = strings.TrimSuffix(subStrMax, regexp.MustCompile(re).FindString(subStrMax))

			// Now, let's apply euristic rules to define what to return...
			// ...
			// If subStrMin=="", output is subStrMax (so, min one word)
			// in other case, return the closest to target ([]runes)
			newLine, excess = ClosestNotEmptyString(target, subStrMin, subStrMax)
			// Food for thoughts:
			//   - prioritize if the subStr ends with a punct
		}
		// Now, update the output and data
		this.translatedLine[i] = newLine
		data = strings.TrimSpace(strings.TrimPrefix(data, newLine))
	}
}
