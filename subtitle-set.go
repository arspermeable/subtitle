package subtitle

import (
	"bufio"
	"io"
)

// -----------------------------------------------
// Functions to import and process the SRT and TRT
// -----------------------------------------------

// Import an SRT file, creating the subtitle blocks and the original text lines
func (this *SubtitleSRT) SetOriginalSrt(reader io.Reader) {
	// Scan the subtitle file for subtitle blocks
	scanner := bufio.NewScanner(reader)
	scanner.Split(SplitSubtitles)
	// Scan and append
	for scanner.Scan() {
		this.appendSubtitle(scanner.Text())
	}
	// Create the slice and underlying array []translatedLine
	this.translatedLine = make([]string, len(this.originalLine))
}

// Import the translated text, into the translatedText field
func (this *SubtitleSRT) SetTranslatedText(txt string) {
	this.translatedText = prepareString(txt)
	this.splitTranslatedTextIntoLineSets()
	for i := range this.lineSet {
		this.splitTranslatedLineSetIntoLines(i)
	}
}

// Import the translated text of a LineSet into its translatedSet field
func (this *SubtitleSRT) SetTranslatedTextOfLineSet(lineSetNumber int, txt string) {
	// Check that lineSet is in range
	if lineSetNumber < 0 || lineSetNumber >= len(this.lineSet) {
		// (****) Should raise an error
		return
	}
	// Assign the txt to the translatedSet
	this.translatedSet[lineSetNumber] = txt
	// Split the translation of this lineSetNumber into lines
	this.splitTranslatedLineSetIntoLines(lineSetNumber)
	// build the translatedText with the new translatedSet
	this.translatedText = joinStrings(this.translatedSet...)
}

// DeleteAllData resets to zero all SubtitleSRT object
func (this *SubtitleSRT) DeleteSubtitleSrt() {
	this.subtitleBlock = nil
	this.lineSet = nil
	this.originalLine = nil
	this.translatedLine = nil
	this.translatedSet = nil
	this.translatedText = ""
}
