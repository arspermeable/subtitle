package subtitle

// ----------------------------------------
// Extract different data from SubtitleSRT
// ----------------------------------------

// Get the original text, returns text and length (runes)
//
func (this *SubtitleSRT) GetOriginalTextFromLines() (string, int) {
	theText := PrepareString(JoinAllStrings(this.originalLine...))
	return theText, len([]rune(theText))
}

// Get the translated text, returns text and length (runes)
//
func (this *SubtitleSRT) GetTranslatedText() (string, int) {
	theText := this.translatedText
	return theText, len([]rune(theText))
}

// Get the translated text using translated lines
// Returns text and length (runes)
func (this *SubtitleSRT) GetTranslatedTextFromLines() (string, int) {
	theText := PrepareString(JoinAllStrings(this.translatedLine...))
	return theText, len([]rune(theText))
}

// Get the translated text using line sets
// Returns text and length (runes)
func (this *SubtitleSRT) GetTranslatedTextFromLineSet() (string, int) {
	theText := PrepareString(JoinAllStrings(this.translatedSet...))
	return theText, len([]rune(theText))
}

// Get the original text of a given line set
// Returns text and length (runes)
func (this *SubtitleSRT) GetOriginalTextOfLineSet(ls int) (string, int) {
	if ls >= len(this.lineSet) || ls < 0 {
		return "", -1
	}
	theText := PrepareString(JoinAllStrings(this.originalLine[this.lineSet[ls].InitLine : this.lineSet[ls].LastLine+1]...))
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

// Get the translated text of a given line set using translated lines
// Returns text and length (runes)
func (this *SubtitleSRT) GetTranslatedTextOfLineSetFromLines(ls int) (string, int) {
	if ls >= len(this.lineSet) || ls < 0 {
		return "", -1
	}
	theText := PrepareString(JoinAllStrings(this.translatedLine[this.lineSet[ls].InitLine : this.lineSet[ls].LastLine+1]...))
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
func (this *SubtitleSRT) GetLineSet() []LineSet {
	return this.lineSet
}
