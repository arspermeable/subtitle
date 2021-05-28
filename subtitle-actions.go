package subtitle

import (
	"fmt"
	"regexp"
	"strings"
)

// --------------------------------------------------
// Functions to manualy manipulate lineSets and lines
// --------------------------------------------------

// MoveLinesFromLineSetToPrev moves n translatedLine(s)
// from top of lineSet to bottom of previous.
// Then, both linesets are processed with SplitTranslatedLineSetIntoLines.
// Input: lineset to move from and number of lines
func (this *SubtitleSRT) MoveLinesFromLineSetToPrev(lsFrom, n int) {
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
	toBeRemoved := joinStrings(this.translatedLine[initLine:endLine]...)

	// remove the text from the top of lineSet[lsFrom].translatedText
	this.translatedSet[lsFrom] = strings.TrimSpace(strings.TrimPrefix(this.translatedSet[lsFrom], toBeRemoved))
	// add the text to the bottom of lineSet[lsFrom-1].translatedText
	this.translatedSet[lsTo] = (this.translatedSet[lsTo] + " " + toBeRemoved)

	// Split the translation of the two affected lineSet into lines
	this.splitTranslatedLineSetIntoLines(lsFrom)
	this.splitTranslatedLineSetIntoLines(lsTo)
}

// MoveWordsFromLineSetToPrev moves n translated word(s)
// from top of lineSet to bottom of previous.
// Then, both linesets are processed with SplitTranslatedLineSetIntoLines.
// Input: lineset to move from and number of words
func (this *SubtitleSRT) MoveWordsFromLineSetToPrev(lsFrom, n int) {
	// Verify that lsFrom is a valid lineset (0 .. #lineSet)
	if lsFrom <= 0 || lsFrom > len(this.lineSet)-1 || n <= 0 {
		return
	}
	// cap n to the number of words
	maxWords := this.CountTranslatedWordsInLineSet(lsFrom)
	if n > maxWords {
		n = maxWords
	}
	// Prepare local variables
	lsTo := lsFrom - 1
	rs := fmt.Sprintf(`^(\S+\s*){%d}`, n)
	loc := regexp.MustCompile((rs)).FindStringIndex(this.translatedSet[lsFrom])
	if loc == nil {
		return
	}

	// Remove the first n words from lsFrom and add it to lsTo
	this.translatedSet[lsTo] = this.translatedSet[lsTo] + " " + strings.TrimSpace(this.translatedSet[lsFrom][loc[0]:loc[1]])
	this.translatedSet[lsFrom] = this.translatedSet[lsFrom][loc[1]:]

	// Split the translation of the two affected lineSet into lines
	this.splitTranslatedLineSetIntoLines(lsFrom)
	this.splitTranslatedLineSetIntoLines(lsTo)
}

// MoveLinesFromLineSetToNext moves n translatedLine(s)
// from bottom of lineSet to top of previous.
// Then, both linesets are processed with SplitTranslatedLineSetIntoLines.
// Input: lineset to move from and number of lines
func (this *SubtitleSRT) MoveLinesFromLineSetToNext(lsFrom, n int) {
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
	toBeRemoved := joinStrings(this.translatedLine[initLine:endLine]...)

	// remove the text from the bottom of lineSet[lsFrom].translatedText
	this.translatedSet[lsFrom] = strings.TrimSpace(strings.TrimSuffix(this.translatedSet[lsFrom], toBeRemoved))
	// add the text to the top of lineSet[lsFrom-1].translatedText
	this.translatedSet[lsTo] = (toBeRemoved + " " + this.translatedSet[lsTo])

	// Split the translation of the two affected lineSet into lines
	this.splitTranslatedLineSetIntoLines(lsFrom)
	this.splitTranslatedLineSetIntoLines(lsTo)
}

// MoveWordsFromLineSetToNext moves n translated words(s)
// from bottom of lineSet to top of previous.
// Then, both linesets are processed with SplitTranslatedLineSetIntoLines.
// Input: lineset to move from and number of words
func (this *SubtitleSRT) MoveWordsFromLineSetToNext(lsFrom, n int) {
	// Verify that lsFrom is a valid lineset (0 .. #lineSet)
	if lsFrom < 0 || lsFrom >= len(this.lineSet)-1 || n <= 0 {
		return
	}
	// cap n to the number of words
	maxWords := this.CountTranslatedWordsInLineSet(lsFrom)
	if n > maxWords {
		n = maxWords
	}
	// Prepare local variables
	lsTo := lsFrom + 1
	rs := fmt.Sprintf(`(\s*\S+){%d}$`, n)
	loc := regexp.MustCompile((rs)).FindStringIndex(this.translatedSet[lsFrom])
	if loc == nil {
		return
	}

	// Remove the last n words from lsFrom and add it to beginning of lsTo
	this.translatedSet[lsTo] = strings.TrimSpace(this.translatedSet[lsFrom][loc[0]:loc[1]] + " " + this.translatedSet[lsTo])
	this.translatedSet[lsFrom] = this.translatedSet[lsFrom][:loc[0]]

	// Split the translation of the two affected lineSet into lines
	this.splitTranslatedLineSetIntoLines(lsFrom)
	this.splitTranslatedLineSetIntoLines(lsTo)
}

// MoveWordFromLineToPrev moves 1 translated words(s)
// from the beginning of a line to the end of the previous one.
// It cannot be used to move words between linesets.
// Affected lineset is *not* processed with SplitTranslatedLineSetIntoLines.
// Input: line number to move from
func (this *SubtitleSRT) MoveWordFromLineToPrev(lineFrom int) {
	// Verify that lineFrom is the first one of a lineset
	if this.IsFirstLineOfLineSet(lineFrom) {
		// if so, fallback to MoveWordsFromLinesetToPrev(1)
		this.MoveWordsFromLineSetToPrev(this.WhatLineSetIsLine(lineFrom), 1)
		return
	}

	// Find the first word of the line
	rs := `^(\S+\s*)`
	loc := regexp.MustCompile((rs)).FindStringIndex(this.translatedLine[lineFrom])
	if loc == nil {
		return
	}
	lineTo := lineFrom - 1

	// Add the first word of lineFrom to lineTo
	this.translatedLine[lineTo] = strings.TrimSpace(this.translatedLine[lineTo] + " " + this.translatedLine[lineFrom][loc[0]:loc[1]])
	// Remove it from lineFrom
	this.translatedLine[lineFrom] = this.translatedLine[lineFrom][loc[1]:]
}

// MoveWordFromLineToNext moves 1 translated words(s)
// from the end of a line to the beginning of the next one.
// It cannot be used to move words between linesets.
// Affected lineset is *not* processed with SplitTranslatedLineSetIntoLines.
// Input: line number to move from
func (this *SubtitleSRT) MoveWordFromLineToNext(lineFrom int) {
	// Verify that lineFrom is the last one of a lineset
	if this.IsLastLineOfLineSet(lineFrom) {
		// if so, fallback to MoveWordsFromLinesetToNext(1)
		this.MoveWordsFromLineSetToNext(this.WhatLineSetIsLine(lineFrom), 1)
		return
	}

	// Find the last word of the line
	rs := `(\s*\S+)$`
	loc := regexp.MustCompile((rs)).FindStringIndex(this.translatedLine[lineFrom])
	if loc == nil {
		return
	}
	lineTo := lineFrom + 1

	// Add the last word of lineFrom to lineTo
	this.translatedLine[lineTo] = strings.TrimSpace(this.translatedLine[lineFrom][loc[0]:loc[1]] + " " + this.translatedLine[lineTo])
	// Remove it from lineFrom
	this.translatedLine[lineFrom] = this.translatedLine[lineFrom][:loc[0]]
}

// SplitLineSetByLine splits a lineset in two by a specified line.
// The break line will be part of the second lineset.
// Then, both linesets are processed with SplitTranslatedLineSetIntoLines.
// It takes the lineset to be split and the line to break at.
func (this *SubtitleSRT) SplitLineSetByLine(ls, breakLine int) {
	// Verify that lsFrom is a valid lineset (0 .. #lineSet)
	if ls < 0 || ls >= len(this.lineSet) {
		return
	}
	// Verify that the lineSet has more than one line, and breakline is in it
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
	this.translatedSet[ls] = joinStrings(this.translatedLine[initLine:breakLine]...)
	this.translatedSet[ls+1] = joinStrings(this.translatedLine[breakLine : lastLine+1]...)
	// *DO NOT* Split again the affected lineSets
	// 20210528: SplitLineSet only splits, it does not re-shufle
	// this.splitTranslatedLineSetIntoLines(ls)
	// this.splitTranslatedLineSetIntoLines(ls + 1)
}

// MergeLineSet merges a LineSet with the previous one
// into a single lineSet
// Resultant lineset is *not* processed with SplitTranslatedLineSetIntoLines.
// It takes the lineset to merge.
func (this *SubtitleSRT) MergeLineSetWithPrev(ls int) {
	// Verify that the situation is legal
	if ls <= 0 || ls > len(this.lineSet)-1 {
		// (****) raise error
		return
	}
	// Last line of LineSet ls-1 now is the last line of ls
	this.lineSet[ls-1].LastLine = this.lineSet[ls].LastLine
	// The translated text of the joint is the joint of the two translated texts
	this.translatedSet[ls-1] = joinStrings(this.translatedSet[ls-1 : ls+1]...)
	// Copy all the subsequent linesets to -1
	this.lineSet = append(this.lineSet[:ls], this.lineSet[ls+1:]...)
	this.translatedSet = append(this.translatedSet[:ls], this.translatedSet[ls+1:]...)
	// 20210528: MergeLineSetWithPrev only merges, it does not split
}

// MergeLineSetWithNext merges a LineSet with the next one
// into a single lineSet
// Resultant lineset is *not* processed with SplitTranslatedLineSetIntoLines.
// It takes the lineset to merge.
func (this *SubtitleSRT) MergeLineSetWithNext(ls int) {
	// Verify that the situation is legal
	if ls < 0 || ls >= len(this.lineSet)-1 {
		// (****) raise error
		return
	}
	// Last line of LineSet ls now is the last line of ls+1
	this.lineSet[ls].LastLine = this.lineSet[ls+1].LastLine
	// The translated text of the joint is the joint of the two translated texts
	this.translatedSet[ls] = joinStrings(this.translatedSet[ls : ls+2]...)
	// Copy all the subsequent linesets to -1
	this.lineSet = append(this.lineSet[:ls+1], this.lineSet[ls+2:]...)
	this.translatedSet = append(this.translatedSet[:ls+1], this.translatedSet[ls+2:]...)
	// 20210528: MergeLineSetWithNext only merges, it does not split
}
