package text

import (
	"fmt"
	"io"
	"time"
)

type consoleWriter struct {
	writer io.Writer

	finalColor     string
	candidateColor string
	stopColor      string
}

func NewConsoleWriter(writer io.Writer) Writer {
	const green = "\033[0;32m"
	const red = "\033[0;31m"
	const stopColor = "\033[K"

	return &consoleWriter{writer: writer, finalColor: green, candidateColor: red, stopColor: stopColor}
}

func (cw *consoleWriter) PrintFinal(_ time.Duration, text string) error {
	// endTime := duration.Truncate(time.Second).String()

	_, err := fmt.Fprintf(cw.writer, "%s%s%s\n", cw.finalColor, text, cw.stopColor)
	if err != nil {
		return fmt.Errorf("cannot print a final message: %w", err)
	}

	return nil
}

func (cw *consoleWriter) PrintCandidate(_ time.Duration, text string) error {
	_, err := fmt.Fprintf(cw.writer, "%s%s%s\r", cw.candidateColor, text, cw.stopColor)
	if err != nil {
		return fmt.Errorf("cannot print a candidate message: %w", err)
	}

	return nil
}

func (cw *consoleWriter) Finalize() error {
	// nothing
	return nil
}
