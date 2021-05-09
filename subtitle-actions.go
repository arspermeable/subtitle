package subtitle

import (
	"fmt"
	"regexp"
	"strings"
)

// --------------------------------------------------
// Functions to manualy manipulate lineSets and lines
// --------------------------------------------------

// Move n translatedLine(s) from top of lineSet to bottom of previous
func (this *SubtitleSRT) MoveLinesUpFromLineSet(lsFrom, n int) {
	// Verify that lsFrom is a valid lineset (1 .. #lineSet-1)
	if lsFrom <= 0 || lsFrom >= len(this.lineSet) || n <= 0 {
		return
	}
	// cap n to the number of lines
	if n > (this.lineSet[lsFrom].LastLine - this.lineSet[lsFrom].InitLine + 1) {
		n = this.lineSet[lsFrom].LastLine - this.lineSet[lsFrom].InitLine + 1
	}
	// Prepare local variables
	lsTo := lsFrom - 1
	initLine := this.lineSet[lsFrom].InitLine
	endLine := this.lineSet[lsFrom].InitLine + n

	// Take last n lines of translastedLine
	toBeRemoved := JoinAllStrings(this.translatedLine[initLine:endLine])

	// remove the text from the top of lineSet[lsFrom].translatedText
	this.translatedSet[lsFrom] = strings.TrimSpace(strings.TrimPrefix(this.translatedSet[lsFrom], toBeRemoved))
	// add the text to the bottom of lineSet[lsFrom-1].translatedText
	this.translatedSet[lsTo] = (this.translatedSet[lsTo] + " " + toBeRemoved)

	// Split the translation of the two affected lineSet into lines
	this.SplitTranslatedLineSetIntoLines(lsFrom)
	this.SplitTranslatedLineSetIntoLines(lsTo)
}

// Move n translated words(s) from top of lineSet to bottom of previous
func (this *SubtitleSRT) MoveWordsUpFromLineSet(lsFrom, n int) {
	// Verify that lsFrom is a valid lineset (0 .. #lineSet)
	if lsFrom <= 0 || lsFrom >= len(this.lineSet) || n <= 0 {
		return
	}
	// cap n to the number of words
	if n > this.CountTranslatedWordsInLineSet(lsFrom) {
		n = this.CountTranslatedWordsInLineSet(lsFrom)
	}
	// Prepare local variables
	lsTo := lsFrom - 1
	rs := fmt.Sprintf(`^(\S+\s+){%d}`, n)
	loc := regexp.MustCompile((rs)).FindStringIndex(this.translatedSet[lsFrom])
	if loc == nil {
		return
	}

	// Remove the first n words from lsFrom and add it to lsTo
	this.translatedSet[lsTo] = this.translatedSet[lsTo] + " " + this.translatedSet[lsFrom][loc[0]:loc[1]]
	this.translatedSet[lsFrom] = this.translatedSet[lsFrom][loc[1]:]

	// Split the translation of the two affected lineSet into lines
	this.SplitTranslatedLineSetIntoLines(lsFrom)
	this.SplitTranslatedLineSetIntoLines(lsTo)
}

// Move words from one line to previous
func (this *SubtitleSRT) MoveWordsUpFromLine(lsFrom, n int) {
}

// Move n translatedLine(s) from bottom of lineSet to top of next
func (this *SubtitleSRT) MoveLinesDownFromLineSet(lsFrom, n int) {
	// Verify that lsFrom is a valid lineset (1 .. #lineSet-1)
	if lsFrom < 0 || lsFrom >= (len(this.lineSet)-1) || n <= 0 {
		return
	}
	// cap n to the number of lines
	if n > (this.lineSet[lsFrom].LastLine - this.lineSet[lsFrom].InitLine + 1) {
		n = this.lineSet[lsFrom].LastLine - this.lineSet[lsFrom].InitLine + 1
	}
	// Prepare local variables
	lsTo := lsFrom + 1
	endLine := this.lineSet[lsFrom].LastLine + 1
	initLine := this.lineSet[lsFrom].LastLine - n + 1

	// Take last n lines of translastedLine
	toBeRemoved := JoinAllStrings(this.translatedLine[initLine:endLine])

	// remove the text from the bottom of lineSet[lsFrom].translatedText
	this.translatedSet[lsFrom] = strings.TrimSpace(strings.TrimSuffix(this.translatedSet[lsFrom], toBeRemoved))
	// add the text to the top of lineSet[lsFrom-1].translatedText
	this.translatedSet[lsTo] = (toBeRemoved + " " + this.translatedSet[lsTo])

	// Split the translation of the two affected lineSet into lines
	this.SplitTranslatedLineSetIntoLines(lsFrom)
	this.SplitTranslatedLineSetIntoLines(lsTo)
}

// Move n translated words(s) from bottom of lineSet to top of previous
func (this *SubtitleSRT) MoveWordsDownFromLineSet(lsFrom, n int) {
	// Verify that lsFrom is a valid lineset (0 .. #lineSet)
	if lsFrom <= 0 || lsFrom >= len(this.lineSet) || n <= 0 {
		return
	}
	// cap n to the number of words
	if n > this.CountTranslatedWordsInLineSet(lsFrom) {
		n = this.CountTranslatedWordsInLineSet(lsFrom)
	}
	// Prepare local variables
	lsTo := lsFrom + 1
	rs := fmt.Sprintf(`(\s+\S+){%d}$`, n)
	loc := regexp.MustCompile((rs)).FindStringIndex(this.translatedSet[lsFrom])
	if loc == nil {
		return
	}

	// Remove the last n words from lsFrom and add it to beginning of lsTo
	this.translatedSet[lsTo] = strings.TrimSpace(this.translatedSet[lsFrom][loc[0]:loc[1]] + " " + this.translatedSet[lsTo])
	this.translatedSet[lsFrom] = this.translatedSet[lsFrom][:loc[0]]

	// Split the translation of the two affected lineSet into lines
	this.SplitTranslatedLineSetIntoLines(lsFrom)
	this.SplitTranslatedLineSetIntoLines(lsTo)
}

// (****)
func (this *SubtitleSRT) MoveWordsDownFromLine(lsFrom, n int) {
}

// Split lineSet ls in two (ls and ls+1) at the originalLine ol
// lineSet[ls] = initLine..ol-1
// lineset[ls+1] = ok..LastLine
func (this *SubtitleSRT) SplitLineSetByLine(ls, breakLine int) {
	// Verify that lsFrom is a valid lineset (0 .. #lineSet)
	if ls < 0 || ls >= len(this.lineSet) {
		return
	}
	// Verify that the lineSet has more than one line, and ol is in it
	numLines := this.lineSet[ls].LastLine - this.lineSet[ls].InitLine + 1
	initLine := this.lineSet[ls].InitLine
	lastLine := this.lineSet[ls].LastLine
	if numLines <= 1 || breakLine <= initLine || breakLine > lastLine {
		return
	}
	// add a new lineSet and translatedSet
	this.lineSet = append(this.lineSet, LineSet{0, 0})
	this.translatedSet = append(this.translatedSet, "")
	// Move lineSets and translatedSet +1
	copy(this.lineSet[ls+1:], this.lineSet[ls:])
	copy(this.translatedSet[ls+1:], this.translatedSet[ls:])
	// adapt the lineSets and translatedSet
	this.lineSet[ls+1].InitLine = breakLine
	this.lineSet[ls+1].LastLine = lastLine
	this.lineSet[ls].LastLine = breakLine - 1
	// Assign the translatedSet
	this.translatedSet[ls] = JoinAllStrings(this.translatedLine[initLine:breakLine])
	this.translatedSet[ls+1] = JoinAllStrings(this.translatedLine[breakLine : lastLine+1])
	// Split again the affected lineSets
	this.SplitTranslatedLineSetIntoLines(ls)
	this.SplitTranslatedLineSetIntoLines(ls + 1)
}

// Merge lineSet ls and ls+1 into a single lineSet
// (****)
func (this *SubtitleSRT) MergeLineSetWithNext(ls int) {
}
