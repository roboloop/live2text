package components

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/roboloop/live2text/internal/services/recognition/text"
)

//go:generate minimock -g -i OutputComponent -s _mock.go -o .

type OutputComponent interface {
	Print(ctx context.Context, writer text.Writer, inputCh <-chan Recognized) error
	ToFile(ctx context.Context, inputCh <-chan Recognized) error
	ToConsole(ctx context.Context, inputCh <-chan Recognized) error
}

type outputComponent struct {
	logger        *slog.Logger
	outputDir     string
	consoleWriter io.Writer
}

func NewOutputComponent(logger *slog.Logger, outputDir string, consoleWriter io.Writer) OutputComponent {
	return &outputComponent{
		logger:        logger,
		outputDir:     outputDir,
		consoleWriter: consoleWriter,
	}
}

func (oc *outputComponent) Print(_ context.Context, writer text.Writer, inputCh <-chan Recognized) error {
	var err error
	for r := range inputCh {
		if r.IsFinal {
			err = writer.PrintFinal(r.EndTime, r.Transcript)
		} else {
			err = writer.PrintCandidate(r.EndTime, r.Transcript)
		}
		if err != nil {
			return fmt.Errorf("cannot print the transcript: %w", err)
		}
	}

	if err = writer.Finalize(); err != nil {
		return fmt.Errorf("cannot finalize the writer: %w", err)
	}

	return nil
}

func (oc *outputComponent) ToFile(ctx context.Context, inputCh <-chan Recognized) error {
	name := time.Now().Format("2006-01-02_15-04-05.txt")
	filename := path.Join(oc.outputDir, name)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("cannot create the file: %w", err)
	}
	defer func() {
		if errClose := file.Close(); errClose != nil {
			oc.logger.ErrorContext(ctx, "Error during closing the file", "error", errClose, "filename", filename)
		}
	}()

	oc.logger.InfoContext(ctx, "Writing the transcript to the file", "filename", filename)

	if err = oc.Print(ctx, text.NewFileWriter(file), inputCh); err != nil {
		return err
	}

	oc.logger.InfoContext(ctx, "Transcript written to the file", "filename", filename)

	return nil
}

func (oc *outputComponent) ToConsole(ctx context.Context, inputCh <-chan Recognized) error {
	return oc.Print(ctx, text.NewConsoleWriter(oc.consoleWriter), inputCh)
}
