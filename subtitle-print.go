package subtitle

import (
	"fmt"
	"io"
)

// --------------------------------------------
// Different ways to print data to an io.Writer
// --------------------------------------------

// Print subtitle block data, one subtitle per line
func (this *SubtitleSRT) PrintSubtitlesData(f io.Writer) {
	for _, sbt := range this.subtitleBlock {
		sbt.PrintShort(f)
	}
}

// Print the line set definition
func (this *SubtitleSRT) PrintLineSetData(f io.Writer) {
	for i, theLineSet := range this.lineSet {
		fmt.Fprintf(f, "Lineset %3.3d, lines: %4.4d-%4.4d, Words: %4d/%-4d, Chars: %5d/%-5d, Ratio: %6.4f, Txt: |>%s<|\n",
			i, theLineSet.InitLine, theLineSet.LastLine,
			this.CountOriginalWordsInLineSet(i), this.CountTranslatedWordsInLineSet(i),
			this.CountOriginalCharsInLineSet(i), this.CountTranslatedCharsInLineSet(i),
			this.CalculateRatioOfLineSet(i),
			PrintStringMaxWidth(this.translatedSet[i], 50))
	}
}

// Print all the original lines, one per line
func (this *SubtitleSRT) PrintOriginalLines(f io.Writer) {
	for i, l := range this.originalLine {
		fmt.Fprintf(f, "%4.4d|%s\n", i, l)
	}
}

// Print all the translated lines, one per line
func (this *SubtitleSRT) PrintTranslatedLines(f io.Writer) {
	for i, l := range this.translatedLine {
		fmt.Fprintf(f, "%4.4d|%s\n", i, l)
	}
}

// Print all the original and translated lines together
func (this *SubtitleSRT) PrintOriginalAndTranslatedLines(f io.Writer) {
	for i, theLineSet := range this.lineSet {
		for j := theLineSet.InitLine; j <= theLineSet.LastLine; j++ {
			fmt.Fprintf(f, "%3.3d|%4.4d|>%s<|>%s<|\n",
				i, j,
				PrintStringMaxWidth(this.originalLine[j], 55),
				PrintStringMaxWidth(this.translatedLine[j], 55))
		}
	}
}

// Print one block with line split data
func (this *SubtitleSRT) PrintLineSetWithSplitData(f io.Writer, theLineSet int) {

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
		fmt.Fprintf(f, "%3.3d|%4.4d|>%s<|>%s<| %4.4d - %4.4d (%5.2f/%+5.2f)\n",
			theLineSet, i,
			PrintStringMaxWidth(this.originalLine[i], 55),
			PrintStringMaxWidth(this.translatedLine[i], 55),
			lenOrig, lenTran, target, excess)
	}
}

// Print all line sets with line split data
func (this *SubtitleSRT) PrintAllLineSetsWithSplitData(f io.Writer) {
	for i := range this.lineSet {
		this.PrintLineSetWithSplitData(f, i)
	}
}

// Print the SRT file, with the original lines
func (this *SubtitleSRT) PrintOriginalSRT(f io.Writer) {

	// Keep count of the lines
	n := 0

	for _, sbt := range this.subtitleBlock {
		sbt.Print(f)
		for i := 0; i < sbt.Nlines; i++ {
			fmt.Fprintln(f, this.originalLine[n])
			n++
		}
		fmt.Fprintln(f)
		fmt.Fprintln(f)

	}
}

// Print the SRT file, with the translated lines
func (this *SubtitleSRT) PrintTranslatedSRT(f io.Writer) {

	// Keep count of the lines
	n := 0

	for _, sbt := range this.subtitleBlock {
		sbt.Print(f)
		for i := 0; i < sbt.Nlines; i++ {
			fmt.Fprintln(f, this.translatedLine[n])
			n++
		}
		fmt.Fprintln(f)
		fmt.Fprintln(f)
	}
}

func (this *SubtitleBlock) PrintShort(f io.Writer) {
	fmt.Fprintf(f, "%s|%s|%2.2d lines|\n", this.Order, this.Timemark, this.Nlines)
}

func (this *SubtitleBlock) Print(f io.Writer) {
	fmt.Fprintln(f, this.Order)
	fmt.Fprintln(f, this.Timemark)
}
