package subtitle

func (this *SubtitleSRT) isFirstLineOfLineSet(theLine int) bool {
	for _, ls := range this.lineSet {
		if ls.InitLine == theLine {
			return true
		}
	}
	return false
}

func (this *SubtitleSRT) isLastLineOfLineSet(theLine int) bool {
	for _, ls := range this.lineSet {
		if ls.LastLine == theLine {
			return true
		}
	}
	return false
}

func (this *SubtitleSRT) isEqual(other SubtitleSRT) bool {
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
