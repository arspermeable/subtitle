package subtitle

// ----------------------------------------
// Extract different data from SubtitleSRT
// ----------------------------------------

// Get the original text, returns text and length (runes)
//
func (this *SubtitleSRT) GetOriginalText() (string, int) {
	// JoinAllLines add a [] for each blank line
	theText := JoinAllLinesWithBrackets(this.originalLine...)
	return theText, len([]rune(theText))
}

// Get the translated text, returns text and length (runes)
//
func (this *SubtitleSRT) GetTranslatedText() (string, int) {
	theText := this.translatedText
	return theText, len([]rune(theText))
}

// Get the original text of a given line set
// Returns text and length (runes)
func (this *SubtitleSRT) GetOriginalTextOfLineSet(ls int) (string, int) {
	if ls >= len(this.lineSet) || ls < 0 {
		return "", -1
	}
	theText := joinStrings(this.originalLine[this.lineSet[ls].InitLine : this.lineSet[ls].LastLine+1]...)
	return theText, len([]rune(theText))
}

// Get the translated text of a given line set
// Returns text and length (runes)
func (this *SubtitleSRT) GetTranslatedTextOfLineSet(ls int) (string, int) {
	if ls >= len(this.lineSet) || ls < 0 {
		return "", -1
	}
	theText := this.translatedSet[ls]
	return theText, len([]rune(theText))
}

// Get the original lines array
// Returns the array with the original lines
func (this *SubtitleSRT) GetOriginalLines() []string {
	return this.originalLine
}

// Get the translated lines array
// Returns the array with the original lines
func (this *SubtitleSRT) GetTranslatedLines() []string {
	return this.translatedLine
}

// GetLineSet returns the LineSet definition in SubtitleSRT
func (this *SubtitleSRT) GetLineSets() []LineSet {
	return this.lineSet
}
