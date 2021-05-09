package subtitle

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// the min length of a line to be matched in the translation
const minmatch = 15
const sep = " ,::!\\.\\?\\)\\]-"

// A subtitle file contains:
//   * an arrray of subtitles, each with 1..n lines of text
//   * an array of lines in the original language
//   * an array of lines in the translated language
//   * an array of line sets definition
//   * an array of the translated text of the LineSet:s
//
// SubtitleFile is based in the SRT definition, each subtitle block consists of
//   <order>     === (\d{1,n})
//   <timemark>  === (hh:mm:ss,mmm --> hh:mm:ss,mmm)
//   <text line> === 0..n lines of text.
//
type SubtitleFile struct {
	subtitle       []Subtitle
	originalLine   []string
	translatedLine []string
	lineSet        []LineSet
	translatedSet  []string
	translatedText string
}

// A LineSet is a set of lines within the list of subtitle text lines
// Each LineSet is process as a block to match original and translation linebreaks
type LineSet struct {
	InitLine int
	LastLine int
}

// Add a subtitle block to the subtitle file
func (this *SubtitleFile) AppendSubtitle(data string) {

	// Split the data by any newline
	lines := regexp.MustCompile(`\r?\n`).Split(data, -1)

	// Extract the lines
	nLine := 2
	for OK := true; OK; OK = nLine < len(lines) {
		var theLine string
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
	this.subtitle = append(this.subtitle, Subtitle{lines[0], lines[1], nLine - 2})
}

// Import an SRT file, creating the subtitle blocks and the original text lines
func (this *SubtitleFile) ImportOriginalSrt(reader io.Reader) {
	// Scan the subtitle file for subtitle blocks
	scanner := bufio.NewScanner(reader)
	scanner.Split(SplitSubtitles)
	// Scan and append
	for scanner.Scan() {
		this.AppendSubtitle(scanner.Text())
	}
	// Create the slice and underlying array []translatedLine
	this.translatedLine = make([]string, len(this.originalLine))
}

// Import the translated text, into the translatedText field
func (this *SubtitleFile) ImportTranslatedText(txt string) {
	this.translatedText = txt
	this.SplitTranslatedTextIntoLineSets()
	for i := range this.lineSet {
		this.SplitTranslatedLineSetIntoLines(i)
	}
}

// Split what is in translatedText, into line sets detected by
// comparing original lines with translatedText
func (this *SubtitleFile) SplitTranslatedTextIntoLineSets() {

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

// Return the ratio translatedChars/originalChars
func (this *SubtitleFile) CalculateRatioOfLineSet(theLineSet int) float64 {
	var ratio float64
	if this.CountOriginalCharsInLineSet(theLineSet) != 0 {
		trChars := this.CountTranslatedCharsInLineSet(theLineSet)
		orChars := this.CountOriginalCharsInLineSet(theLineSet)
		ratio = float64(trChars) / float64(orChars)
	}
	return ratio
}

// Split the text assigned to a LineSet into lines
func (this *SubtitleFile) SplitTranslatedLineSetIntoLines(theLineSet int) {

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

// ------------------------------------------
// Functions to manipulate lineSets and lines
// ------------------------------------------

// Move n translatedLine(s) from top of lineSet to bottom of previous
func (this *SubtitleFile) MoveLinesUpFromLineSet(lsFrom, n int) {
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
func (this *SubtitleFile) MoveWordsUpFromLineSet(lsFrom, n int) {
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
func (this *SubtitleFile) MoveWordsUpFromLine(lsFrom, n int) {
}

// Move n translatedLine(s) from bottom of lineSet to top of next
func (this *SubtitleFile) MoveLinesDownFromLineSet(lsFrom, n int) {
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
func (this *SubtitleFile) MoveWordsDownFromLineSet(lsFrom, n int) {
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

func (this *SubtitleFile) MoveWordsDownFromLine(lsFrom, n int) {
}

// Split lineSet ls in two (ls and ls+1) at the originalLine ol
// lineSet[ls] = initLine..ol-1
// lineset[ls+1] = ok..LastLine
func (this *SubtitleFile) SplitLineSetByLine(ls, breakLine int) {
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
func (this *SubtitleFile) MergeLineSetWithNext(ls int) {
}

// --------------------------------------------
// Informative functions about the SubtitleFile
// --------------------------------------------

// CountLineSets returns the total number of LineSet:s in SubtitleFile
func (this *SubtitleFile) CountLineSets() int {
	return len(this.lineSet)
}

// CountLines returns the total number of lines in SubtitleFile
func (this *SubtitleFile) CountLines() int {
	return len(this.originalLine)
}

// CountOriginalWords returns the total number of original words in SubtitleFile
func (this *SubtitleFile) CountOriginalWords() int {
	originalText, _ := this.GetOriginalTextFromLines()
	return len(strings.Fields(originalText))
}

// CountOriginalWordsInLine returns the total number of original words in a given line
func (this *SubtitleFile) CountOriginalWordsInLine(theLine int) int {
	if theLine >= len(this.originalLine) || theLine < 0 {
		return -1
	}

	return len(strings.Fields(this.originalLine[theLine]))
}

// CountOriginalChars returns the original chars in SubtitleFile
func (this *SubtitleFile) CountOriginalChars() int {
	originalText, _ := this.GetOriginalTextFromLines()
	return len([]rune(originalText)) - len(this.originalLine) + 1
}

// Count original chars in a given line
func (this *SubtitleFile) CountOriginalCharsInLine(theLine int) int {
	if theLine >= len(this.originalLine) || theLine < 0 {
		return -1
	}

	return len([]rune(this.originalLine[theLine]))
}

// Count lines in a given line set
func (this *SubtitleFile) CountLinesInLineSet(theLineSet int) int {
	if theLineSet >= len(this.lineSet) || theLineSet < 0 {
		return -1
	}

	return this.lineSet[theLineSet].LastLine - this.lineSet[theLineSet].InitLine + 1
}

// Count original words in a given line set
func (this *SubtitleFile) CountOriginalWordsInLineSet(theLineSet int) int {
	if theLineSet >= len(this.lineSet) || theLineSet < 0 {
		return -1
	}

	text, _ := this.GetOriginalTextOfLineSet(theLineSet)
	return len(strings.Fields(text))
}

// Count original chars (runes) in a given line set
func (this *SubtitleFile) CountOriginalCharsInLineSet(theLineSet int) int {
	if theLineSet >= len(this.lineSet) || theLineSet < 0 {
		return -1
	}
	var origChars int
	for i := this.lineSet[theLineSet].InitLine; i <= this.lineSet[theLineSet].LastLine; i++ {
		origChars += this.CountOriginalCharsInLine(i)
	}
	return origChars
}

// CountTranslatedWords returns the number of translated words in a SubtitleFile
func (this *SubtitleFile) CountTranslatedWords() int {
	return len(strings.Fields(this.translatedText))
}

// Count translated words in a given line set
func (this *SubtitleFile) CountTranslatedWordsInLineSet(theLineSet int) int {
	if theLineSet >= len(this.lineSet) || theLineSet < 0 {
		return -1
	}

	return len(strings.Fields(this.translatedSet[theLineSet]))
}

// CountTranslatedChars returns translated chars (runes) in a SubtitleFile
func (this *SubtitleFile) CountTranslatedChars() int {
	// Translated Chars is number of chars - CRLF (number of lines + 1)
	return len([]rune(this.translatedText)) - len(this.originalLine) + 1
}

// Count translated chars (runes) in a given line set
func (this *SubtitleFile) CountTranslatedCharsInLineSet(theLineSet int) int {
	if theLineSet >= len(this.lineSet) || theLineSet < 0 {
		return -1
	}

	// NUmber of Chars is total runes - CRLFs (****) more precise: CRLFs - empty lines
	return len([]rune(this.translatedSet[theLineSet])) - (this.lineSet[theLineSet].LastLine - this.lineSet[theLineSet].InitLine)
}

// ----------------------------------------
// Extract different data from SubtitleFile
// ----------------------------------------

// Get the original text, returns text and length (runes)
//
func (this *SubtitleFile) GetOriginalTextFromLines() (string, int) {
	theText := PrepareString(JoinAllStrings(this.originalLine))
	return theText, len([]rune(theText))
}

// Get the translated text, returns text and length (runes)
//
func (this *SubtitleFile) GetTranslatedText() (string, int) {
	theText := this.translatedText
	return theText, len([]rune(theText))
}

// Get the translated text using translated lines
// Returns text and length (runes)
func (this *SubtitleFile) GetTranslatedTextFromLines() (string, int) {
	theText := PrepareString(JoinAllStrings(this.translatedLine))
	return theText, len([]rune(theText))
}

// Get the translated text using line sets
// Returns text and length (runes)
func (this *SubtitleFile) GetTranslatedTextFromLineSet() (string, int) {
	theText := PrepareString(JoinAllStrings(this.translatedSet))
	return theText, len([]rune(theText))
}

// Get the original text of a given line set
// Returns text and length (runes)
func (this *SubtitleFile) GetOriginalTextOfLineSet(ls int) (string, int) {
	if ls >= len(this.lineSet) || ls < 0 {
		return "", -1
	}
	theText := PrepareString(JoinAllStrings(this.originalLine[this.lineSet[ls].InitLine : this.lineSet[ls].LastLine+1]))
	return theText, len([]rune(theText))
}

// Get the translated text of a given line set
// Returns text and length (runes)
func (this *SubtitleFile) GetTranslatedTextOfLineSet(ls int) (string, int) {
	if ls >= len(this.lineSet) || ls < 0 {
		return "", -1
	}
	theText := this.translatedSet[ls]
	return theText, len([]rune(theText))
}

// Get the translated text of a given line set using translated lines
// Returns text and length (runes)
func (this *SubtitleFile) GetTranslatedTextOfLineSetFromLines(ls int) (string, int) {
	if ls >= len(this.lineSet) || ls < 0 {
		return "", -1
	}
	theText := PrepareString(JoinAllStrings(this.translatedLine[this.lineSet[ls].InitLine : this.lineSet[ls].LastLine+1]))
	return theText, len([]rune(theText))
}

// Get the original lines array
// Returns the array with the original lines
func (this *SubtitleFile) GetOriginalLines() []string {
	return this.originalLine
}

// Get the translated lines array
// Returns the array with the original lines
func (this *SubtitleFile) GetTranslatedLines() []string {
	return this.translatedLine
}

// GetLineSet returns the LineSet definition in SubtitleFile
func (this *SubtitleFile) GetLineSet() []LineSet {
	return this.lineSet
}

// --------------------------------------------
// Different ways to print data to an io.Writer
// --------------------------------------------

// Print subtitle block data, one subtitle per line
func (this *SubtitleFile) PrintSubtitlesData(f io.Writer) {
	for _, sbt := range this.subtitle {
		sbt.PrintShort(f)
	}
}

// Print the line set definition
func (this *SubtitleFile) PrintLineSetData() {
	for i, theLineSet := range this.lineSet {
		fmt.Printf("Lineset %3.3d, lines: %4.4d-%4.4d, Words: %4d/%-4d, Chars: %5d/%-5d, Ratio: %6.4f, Txt: |>%s<|\n",
			i, theLineSet.InitLine, theLineSet.LastLine,
			this.CountOriginalWordsInLineSet(i), this.CountTranslatedWordsInLineSet(i),
			this.CountOriginalCharsInLineSet(i), this.CountTranslatedCharsInLineSet(i),
			this.CalculateRatioOfLineSet(i),
			PrintStringMaxWidth(this.translatedSet[i], 50))
	}
}

// Print all the original lines, one per line
func (this *SubtitleFile) PrintOriginalLines() {
	for i, l := range this.originalLine {
		fmt.Printf("%4.4d|%s\n", i, l)
	}
}

// Print all the translated lines, one per line
func (this *SubtitleFile) PrintTranslatedLines() {
	for i, l := range this.translatedLine {
		fmt.Printf("%4.4d|%s\n", i, l)
	}
}

// Print all the original and translated lines together
func (this *SubtitleFile) PrintOriginalAndTranslatedLines() {
	for i, theLineSet := range this.lineSet {
		for j := theLineSet.InitLine; j <= theLineSet.LastLine; j++ {
			fmt.Printf("%3.3d|%4.4d|>%s<|>%s<|\n",
				i, j,
				PrintStringMaxWidth(this.originalLine[j], 55),
				PrintStringMaxWidth(this.translatedLine[j], 55))
		}
	}
}

// Print one block with line split data
func (this *SubtitleFile) PrintLineSetWithSplitData(theLineSet int) {

	// Calculate the ratio translated:original for this line set
	ratio := 0.0
	if this.CountOriginalCharsInLineSet(theLineSet) != 0 {
		trChars := this.CountTranslatedCharsInLineSet(theLineSet)
		orChars := this.CountOriginalCharsInLineSet(theLineSet)
		ratio = float64(trChars) / float64(orChars)
	}

	var lenOrig, lenTran int
	var target, excess float64

	// Iterate over the lines of the lineSet theLineSet
	for i := this.lineSet[theLineSet].InitLine; i <= this.lineSet[theLineSet].LastLine; i++ {

		lenOrig = len([]rune(this.originalLine[i]))
		lenTran = len([]rune(this.translatedLine[i]))
		target = ratio*float64(lenOrig) - excess
		excess = float64(lenTran) - target

		// Print the text
		fmt.Printf("%3.3d|%4.4d|>%s<|>%s<| %4.4d - %4.4d (%5.2f/%+5.2f)\n",
			theLineSet, i,
			PrintStringMaxWidth(this.originalLine[i], 55),
			PrintStringMaxWidth(this.translatedLine[i], 55),
			lenOrig, lenTran, target, excess)
	}
}

// Print all line sets with line split data
func (this *SubtitleFile) PrintAllLineSetsWithSplitData() {
	for i := range this.lineSet {
		this.PrintLineSetWithSplitData(i)
	}
}

// Print the SRT file, with the original lines
func (this *SubtitleFile) PrintOriginalSRT() {

	// Keep count of the lines
	n := 0

	for _, sbt := range this.subtitle {
		sbt.Print()
		for i := 0; i < sbt.Nlines; i++ {
			fmt.Println(this.originalLine[n])
			n++
		}
		fmt.Println()
		fmt.Println()

	}
}

// Print the SRT file, with the translated lines
func (this *SubtitleFile) PrintTranslatedSRT() {

	// Keep count of the lines
	n := 0

	for _, sbt := range this.subtitle {
		sbt.Print()
		for i := 0; i < sbt.Nlines; i++ {
			fmt.Println(this.translatedLine[n])
			n++
		}
		fmt.Println()
		fmt.Println()

	}
}
