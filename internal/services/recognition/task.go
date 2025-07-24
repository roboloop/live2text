package recognition

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/sync/errgroup"

	"live2text/internal/services/audio"
	"live2text/internal/services/recognition/components"
	"live2text/internal/services/recognition/text"
	"live2text/internal/utils"
)

type Task struct {
	logger    *slog.Logger
	audio     audio.Audio
	burner    components.BurnerComponent
	recognize components.RecognizerComponent
	socket    components.SocketComponent
	output    components.OutputComponent

	formatter  *text.Formatter
	device     string
	language   string
	socketPath string
}

func NewTask(
	logger *slog.Logger,
	audio audio.Audio,
	burner components.BurnerComponent,
	recognize components.RecognizerComponent,
	socket components.SocketComponent,
	output components.OutputComponent,
	device, language, socketPath string,
) *Task {
	return &Task{
		logger:    logger.With("component", "task"),
		audio:     audio,
		burner:    burner,
		recognize: recognize,
		socket:    socket,
		output:    output,

		formatter:  text.NewSubtitleFormatter(0, 0),
		device:     device,
		language:   language,
		socketPath: socketPath,
	}
}

func (t *Task) Run(ctx context.Context) error {
	deviceListener, err := t.audio.ListenDevice(t.device)
	if err != nil {
		return fmt.Errorf("cannot listen to a device: %w", err)
	}

	audioBroadcaster := utils.Broadcaster(ctx, t.logger, deviceListener.GetChannel(), []string{"burner", "recognize"})
	burnerCh, recognizeCh := audioBroadcaster[0], audioBroadcaster[1]

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return deviceListener.Listen(ctx)
	})
	g.Go(func() error {
		parameters := deviceListener.GetParameters()
		return t.burner.SaveAudio(ctx, burnerCh, components.BurnerParameters{
			Channels:   parameters.Channels,
			SampleRate: parameters.SampleRate,
		})
	})

	recognizedCh := make(chan components.Recognized, 100)
	g.Go(func() error {
		defer close(recognizedCh)

		parameters := deviceListener.GetParameters()
		return t.recognize.Recognize(ctx, components.RecognizeParameters{
			Channels:   parameters.Channels,
			SampleRate: parameters.SampleRate,
			Language:   t.language,
		}, recognizeCh, recognizedCh)
	})

	recognizedBroadcaster := utils.Broadcaster(ctx, t.logger, recognizedCh, []string{"console", "file", "subtitle"})
	consoleCh, fileCh, subtitleCh := recognizedBroadcaster[0], recognizedBroadcaster[1], recognizedBroadcaster[2]
	formatter := text.NewSubtitleFormatter(0, 0)
	g.Go(func() error {
		return t.output.ToConsole(ctx, consoleCh)
	})
	g.Go(func() error {
		return t.output.ToFile(ctx, fileCh)
	})
	g.Go(func() error {
		return t.output.Print(ctx, text.NewSubtitleWriter(formatter), subtitleCh)
	})

	g.Go(func() error {
		return t.socket.Listen(ctx, t.socketPath, formatter)
	})

	if err = g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		t.logger.ErrorContext(ctx, "Recognition Task failed", "error", err)
	}

	t.logger.InfoContext(ctx, "Shutting down Task...", "device", t.device)

	return err
}

func (t *Task) Text() string {
	return t.formatter.Format()
}

func (t *Task) SocketPath() string {
	return t.socketPath
}
