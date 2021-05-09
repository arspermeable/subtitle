package subtitle

import (
	"fmt"
	"io"
)

// A subtitle defines a subtitle block in a SRT file
//   * the order (\d{1,n})
//   * the time mark (hh:mm:ss,mmm --> hh:mm:ss,mmm)
//   * Number of lines (1..n).
// Note: If a subtitle includes no lines, then nLines=1 and "" is added
type Subtitle struct {
	Order    string
	Timemark string
	Nlines   int
}

func (this *Subtitle) PrintShort(f io.Writer) {
	fmt.Fprintf(f, "%s|%s|%2.2d lines|\n", this.Order, this.Timemark, this.Nlines)
}

func (this *Subtitle) Print(f io.Writer) {
	fmt.Fprintln(f, this.Order)
	fmt.Fprintln(f, this.Timemark)
}
