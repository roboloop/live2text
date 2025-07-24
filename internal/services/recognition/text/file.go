package text

import (
	"fmt"
	"io"
	"time"
)

type fileWriter struct {
	writer io.Writer

	latestCandidate string
}

func NewFileWriter(writer io.Writer) Writer {
	return &fileWriter{writer: writer}
}

func (fw *fileWriter) PrintFinal(_ time.Duration, text string) error {
	_, err := fmt.Fprintf(fw.writer, "%s\n", text)
	if err != nil {
		return fmt.Errorf("cannot print a final message: %w", err)
	}

	fw.latestCandidate = ""

	return nil
}

func (fw *fileWriter) PrintCandidate(_ time.Duration, text string) error {
	fw.latestCandidate = text
	return nil
}

func (fw *fileWriter) Finalize() error {
	if fw.latestCandidate == "" {
		return nil
	}

	return fw.PrintFinal(0, fw.latestCandidate)
}
