package components_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/services/burner"
	"github.com/roboloop/live2text/internal/services/recognition/components"
	"github.com/roboloop/live2text/internal/utils/logger"
)

func TestBurner(t *testing.T) {
	t.Parallel()

	t.Run("cannot create the file", func(t *testing.T) {
		t.Parallel()

		bc := setupBurner(t, nil, nil, string([]byte{0x00}))

		err := bc.SaveAudio(t.Context(), nil, components.BurnerParameters{})

		require.Error(t, err)
		require.ErrorContains(t, err, "cannot create the file")
	})

	t.Run("cannot burn the file", func(t *testing.T) {
		t.Parallel()

		bc := setupBurner(t, func(mc *minimock.Controller, b *burner.BurnerMock) {
			b.BurnMock.Return(errors.New("dummy error"))
		}, nil, "")

		err := bc.SaveAudio(t.Context(), nil, components.BurnerParameters{})

		require.Error(t, err)
		require.ErrorContains(t, err, "dummy error")
		require.ErrorContains(t, err, "cannot burn the file")
	})

	t.Run("burns a file to disk", func(t *testing.T) {
		t.Parallel()

		ch := make(chan []int16)
		defer close(ch)

		bc := setupBurner(t, func(mc *minimock.Controller, b *burner.BurnerMock) {
			b.BurnMock.
				Set(
					func(ctx context.Context, w io.Writer, input <-chan []int16, channels, sampleRate int) error {
						var inputCh <-chan []int16 = ch
						require.Equal(t, inputCh, input)
						require.Equal(t, 1, channels)
						require.Equal(t, 24000, sampleRate)

						<-ctx.Done()
						return nil
					},
				)
		}, logger.NilLogger, "")
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		err := bc.SaveAudio(ctx, ch, components.BurnerParameters{
			Channels:   1,
			SampleRate: 24000,
		})

		require.NoError(t, err)
	})
}

func setupBurner(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, b *burner.BurnerMock),
	l *slog.Logger,
	outputDir string,
) components.BurnerComponent {
	mc := minimock.NewController(t)
	b := burner.NewBurnerMock(mc)

	if setupMocks != nil {
		setupMocks(mc, b)
	}
	if l == nil {
		l = logger.NilLogger
	}
	if outputDir == "" {
		outputDir = t.TempDir()
	}

	return components.NewBurnerComponent(l, b, outputDir)
}
