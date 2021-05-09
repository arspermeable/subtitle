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
	order    string
	timemark string
	nLines   int
}

func (this *Subtitle) PrintShort(f io.Writer) {
	fmt.Fprintf(f, "%s|%s|%2.2d lines|\n", this.order, this.timemark, this.nLines)
}

func (this *Subtitle) Print() {
	fmt.Println(this.order)
	fmt.Println(this.timemark)
}
