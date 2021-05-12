package subtitle

import "strings"

// --------------------------------------------
// Informative functions about the SubtitleSRT
// --------------------------------------------

// Return the ratio translatedChars/originalChars
func (this *SubtitleSRT) CalculateRatioOfLineSet(theLineSet int) float64 {
	var ratio float64
	if this.CountOriginalCharsInLineSet(theLineSet) != 0 {
		trChars := this.CountTranslatedCharsInLineSet(theLineSet)
		orChars := this.CountOriginalCharsInLineSet(theLineSet)
		ratio = float64(trChars) / float64(orChars)
	}
	return ratio
}

// CountLineSets returns the total number of LineSet:s in SubtitleSRT
func (this *SubtitleSRT) CountLineSets() int {
	return len(this.lineSet)
}

// CountLines returns the total number of lines in SubtitleSRT
func (this *SubtitleSRT) CountLines() int {
	return len(this.originalLine)
}

// CountOriginalWords returns the total number of original words in SubtitleSRT
func (this *SubtitleSRT) CountOriginalWords() int {
	originalText, _ := this.GetOriginalTextFromLines()
	return len(strings.Fields(originalText))
}

// Count lines in a given line set
func (this *SubtitleSRT) CountLinesInLineSet(theLineSet int) int {
	if theLineSet >= len(this.lineSet) || theLineSet < 0 {
		return -1
	}

	return this.lineSet[theLineSet].LastLine - this.lineSet[theLineSet].InitLine + 1
}

// CountOriginalWordsInLine returns the total number of original words in a given line
func (this *SubtitleSRT) CountOriginalWordsInLine(theLine int) int {
	if theLine >= len(this.originalLine) || theLine < 0 {
		return -1
	}

	return len(strings.Fields(this.originalLine[theLine]))
}

// CountOriginalChars returns the original chars in SubtitleSRT
func (this *SubtitleSRT) CountOriginalChars() int {
	originalText, _ := this.GetOriginalTextFromLines()
	return len([]rune(originalText)) - len(this.originalLine) + 1
}

// Count original chars in a given line
func (this *SubtitleSRT) CountOriginalCharsInLine(theLine int) int {
	if theLine >= len(this.originalLine) || theLine < 0 {
		return -1
	}

	return len([]rune(this.originalLine[theLine]))
}

// Count original words in a given line set
func (this *SubtitleSRT) CountOriginalWordsInLineSet(theLineSet int) int {
	if theLineSet >= len(this.lineSet) || theLineSet < 0 {
		return -1
	}

	text, _ := this.GetOriginalTextOfLineSet(theLineSet)
	return len(strings.Fields(text))
}

// Count original chars (runes) in a given line set
func (this *SubtitleSRT) CountOriginalCharsInLineSet(theLineSet int) int {
	if theLineSet >= len(this.lineSet) || theLineSet < 0 {
		return -1
	}
	var origChars int
	for i := this.lineSet[theLineSet].InitLine; i <= this.lineSet[theLineSet].LastLine; i++ {
		origChars += this.CountOriginalCharsInLine(i)
	}
	return origChars
}

// CountTranslatedWords returns the number of translated words in a SubtitleSRT
func (this *SubtitleSRT) CountTranslatedWords() int {
	return len(strings.Fields(this.translatedText))
}

// Count translated words in a given line set
func (this *SubtitleSRT) CountTranslatedWordsInLineSet(theLineSet int) int {
	if theLineSet >= len(this.lineSet) || theLineSet < 0 {
		return -1
	}

	return len(strings.Fields(this.translatedSet[theLineSet]))
}

// CountTranslatedChars returns translated chars (runes) in a SubtitleSRT
func (this *SubtitleSRT) CountTranslatedChars() int {
	// Translated Chars is number of chars - CRLF (number of lines + 1)
	return len([]rune(this.translatedText)) - len(this.originalLine) + 1
}

// Count translated chars (runes) in a given line set
func (this *SubtitleSRT) CountTranslatedCharsInLineSet(theLineSet int) int {
	if theLineSet >= len(this.lineSet) || theLineSet < 0 {
		return -1
	}

	// NUmber of Chars is total runes - CRLFs (****) more precise: CRLFs - empty lines
	return len([]rune(this.translatedSet[theLineSet])) - (this.lineSet[theLineSet].LastLine - this.lineSet[theLineSet].InitLine)
}
