package components

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"time"

	"live2text/internal/services/burner"
)

type BurnerParameters struct {
	Channels   int
	SampleRate int
}

//go:generate minimock -g -i BurnerComponent -s _mock.go -o .

type BurnerComponent interface {
	SaveAudio(ctx context.Context, input <-chan []int16, parameters BurnerParameters) error
}

type burnerComponent struct {
	logger    *slog.Logger
	burner    burner.Burner
	outputDir string
}

func NewBurnerComponent(logger *slog.Logger, burner burner.Burner, outputDir string) BurnerComponent {
	return &burnerComponent{logger: logger, burner: burner, outputDir: outputDir}
}

func (b *burnerComponent) SaveAudio(ctx context.Context, input <-chan []int16, parameters BurnerParameters) error {
	name := time.Now().Format("2006-01-02_15-04-05.wav")
	filename := path.Join(b.outputDir, name)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("cannot create the file: %w", err)
	}
	defer func() {
		if errClose := file.Close(); errClose != nil {
			b.logger.ErrorContext(ctx, "Error closing the file", "error", errClose, "filename", filename)
		}
	}()

	if err = b.burner.Burn(ctx, file, input, parameters.Channels, parameters.SampleRate); err != nil {
		return fmt.Errorf("cannot burn the file: %w", err)
	}

	return nil
}
