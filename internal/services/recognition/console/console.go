package console

import (
	"fmt"
	"io"
)

type Writer struct {
	w io.Writer
}

const (
	red    = "\033[0;31m"
	green  = "\033[0;32m"
	yellow = "\033[0;33m"

	erase = "\033[K"
)

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

func (w *Writer) PrintSuccess(text string) {
	fmt.Fprintf(w.w, "%s%s%s\n", green, erase, text)
}

func (w *Writer) PrintFail(text string) {
	fmt.Fprintf(w.w, "%s%s%s\r", red, erase, text)
}

func (w *Writer) PrintNewLine() {
	fmt.Fprintf(w.w, "\n")
}
