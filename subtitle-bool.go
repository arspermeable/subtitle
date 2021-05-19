package subtitle

import (
	"regexp"
	"strings"
)

// IsLoaded returns true if the data has been loaded into this
// Data loaded means that the original SRT and the translation TRT have been loaded
// and the text has been split into LineSets
func (this *SubtitleSRT) IsLoaded() bool {
	return this.IsLoadedSRT() && this.IsLoadedTRT() && this.IsSplit()
}

// IsLoadedSRT returns true if the original SRT has been loaded into this
// SRT loaded means that subtitleBlock and originalLine have length > 0
func (this *SubtitleSRT) IsLoadedSRT() bool {
	return len(this.subtitleBlock) > 0 && len(this.originalLine) > 0
}

// IsLoadedTRT returns true if the translated TRT has been loaded into this
// TRT loaded means that translatedText is not an empty string
func (this *SubtitleSRT) IsLoadedTRT() bool {
	return this.translatedText != ""
}

// IsSplit returns true if the text has been split in linesets
// Text split means that lineSet and translatedSet has length>0,
// and translatedLine has length > 0 (although this is trivial per implementation)
func (this *SubtitleSRT) IsSplit() bool {
	return len(this.lineSet) > 0 && len(this.translatedSet) > 0 && len(this.translatedLine) > 0
}

// isFirstLineOfLineSet returns true if the line is the first of any lineSet
func (this *SubtitleSRT) IsFirstLineOfLineSet(theLine int) bool {
	for _, ls := range this.lineSet {
		if ls.InitLine == theLine {
			return true
		}
	}
	return false
}

// isLastLineOfLineSet returns true if the line is the last of any lineSet
func (this *SubtitleSRT) IsLastLineOfLineSet(theLine int) bool {
	for _, ls := range this.lineSet {
		if ls.LastLine == theLine {
			return true
		}
	}
	return false
}

// isEqual returns true if the two SubtitleSRT are the same
// It's a custom deep comparison of struct/slices
func (this *SubtitleSRT) IsEqual(other SubtitleSRT) bool {
	// subtitleBlock
	if len(this.subtitleBlock) != len(other.subtitleBlock) {
		return false
	}
	for i, b := range this.subtitleBlock {
		if b != other.subtitleBlock[i] {
			return false
		}
	}
	// lineSet
	if len(this.lineSet) != len(other.lineSet) {
		return false
	}
	for i, l := range this.lineSet {
		if l != other.lineSet[i] {
			return false
		}
	}
	// originalLine
	if len(this.originalLine) != len(other.originalLine) {
		return false
	}
	for i, o := range this.originalLine {
		if o != other.originalLine[i] {
			return false
		}
	}
	// translatedLine
	if len(this.translatedLine) != len(other.translatedLine) {
		return false
	}
	for i, t := range this.translatedLine {
		if t != other.translatedLine[i] {
			return false
		}
	}
	// translatedSet
	if len(this.translatedSet) != len(other.translatedSet) {
		return false
	}
	for i, t := range this.translatedSet {
		if t != other.translatedSet[i] {
			return false
		}
	}
	// translatedText
	if this.translatedText != other.translatedText {
		return false
	}
	return true
}

// isTranslationConsistent returns true if the translatedText,
// the translatedText from linesets and the translatedText from
// lines is the same.
func (this *SubtitleSRT) IsTranslationConsistent() bool {
	// Verify the translatedLines
	theTextLines := joinStrings(this.translatedLine...)
	theTextSets := joinStrings(this.translatedSet...)
	theText := regexp.MustCompile(`\s*\[\]\s*`).ReplaceAllString(this.translatedText, " ")
	if !strings.EqualFold(theTextLines, theTextSets) {
		return false
	}
	if !strings.EqualFold(theText, theTextLines) {
		return false
	}
	return true
}
