package burner_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"live2text/internal/services/burner"
	"live2text/internal/services/metrics"
	"live2text/internal/utils/logger"
)

type failingWriter struct {
	bytes.Buffer
}

func (fw *failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("written error")
}

func TestBurn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string

		writer   io.Writer
		channels int

		expectedWritten []int16
		expectedErr     string
	}{
		{
			"burn mono audio",
			&bytes.Buffer{},
			1,
			[]int16{},
			"",
		},
		{
			"burn stereo audio",
			&bytes.Buffer{},
			2,
			[]int16{},
			"",
		},
		{
			"burn stereo audio",
			&failingWriter{},
			2,
			[]int16{},
			"written error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			b := burner.NewBurner(logger.NilLogger, metrics.NewMetrics(nil, nil))
			input := make(chan []int16, 10)
			rawInt16 := []int16{255, 256, 257, 258}
			rawUint8 := []uint8{
				255, 0,
				0, 1,
				1, 1,
				2, 1,
			}
			input <- rawInt16
			go func() {
				time.Sleep(10 * time.Millisecond)
				cancel()
			}()

			err := b.Burn(ctx, tt.writer, input, tt.channels, 48000)
			if tt.expectedErr != "" {
				require.ErrorContains(t, err, tt.expectedErr)
				return
			}

			require.NoError(t, err)
			buf := tt.writer.(*bytes.Buffer)
			require.NotZero(t, buf.Len())
			require.Equal(t, rawUint8, buf.Bytes()[buf.Len()-len(rawInt16)*2:])
		})
	}
}
