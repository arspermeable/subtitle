//
// Package subtitle modelizes an SRT file and its translation into a second language.
//
//
// Some basic operations to manage subtitles are added
//
package subtitle

// A subtitle defines a subtitle block in a SRT file
//   * the order (\d{1,n})
//   * the time mark (hh:mm:ss,mmm --> hh:mm:ss,mmm)
//   * Number of lines (1..n).
type SubtitleBlock struct {
	Order    string
	Timemark string
	Nlines   int
}

// A LineSet is a set of lines within the list of subtitle text lines
// Each LineSet is process as a block to match original and translation linebreaks
// (****) This may be simplified if LastLine is not included (as in an array)
type LineSet struct {
	InitLine int
	LastLine int
}

// A subtitle file contains:
//   * an arrray of subtitles, each with 1..n lines of text
//   * an array of lines in the original language
//   * an array of lines in the translated language
//   * an array of line sets definition
//   * an array of the translated text of the LineSet:s
//
// SubtitleSRT is based in the SRT definition, each subtitle block consists of
//   <order>     === (\d{1,n})
//   <timemark>  === (hh:mm:ss,mmm --> hh:mm:ss,mmm)
//   <text line> === 0..n lines of text.
//
type SubtitleSRT struct {
	subtitleBlock  []SubtitleBlock
	lineSet        []LineSet
	originalLine   []string
	translatedLine []string
	translatedSet  []string
	translatedText string
}
