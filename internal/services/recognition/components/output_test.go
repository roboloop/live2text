package components_test

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/services/recognition/components"
	"github.com/roboloop/live2text/internal/services/recognition/text"
	"github.com/roboloop/live2text/internal/utils/logger"
)

func TestPrint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		recognized []components.Recognized
		setupMocks func(mc *minimock.Controller, w *text.WriterMock)

		expectErr string
	}{
		{
			name:       "no data",
			recognized: []components.Recognized{},
			setupMocks: func(mc *minimock.Controller, w *text.WriterMock) {
				w.FinalizeMock.Return(nil)
			},
			expectErr: "",
		},
		{
			name: "cannot print a final transcript",
			recognized: []components.Recognized{{
				Transcript: "foo",
				IsFinal:    true,
			}},
			setupMocks: func(mc *minimock.Controller, w *text.WriterMock) {
				w.PrintFinalMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot print the transcript",
		},
		{
			name: "cannot print a candidate transcript",
			recognized: []components.Recognized{{
				Transcript: "foo",
				IsFinal:    false,
			}},
			setupMocks: func(mc *minimock.Controller, w *text.WriterMock) {
				w.PrintCandidateMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot print the transcript",
		},
		{
			name:       "cannot finalize the writer",
			recognized: []components.Recognized{},
			setupMocks: func(mc *minimock.Controller, w *text.WriterMock) {
				w.FinalizeMock.Return(errors.New("something happened"))
			},
			expectErr: "cannot finalize the writer",
		},
		{
			name: "print recognized data",
			recognized: []components.Recognized{{
				Transcript: "foo",
				IsFinal:    false,
				EndTime:    100 * time.Millisecond,
			}, {
				Transcript: "bar",
				IsFinal:    true,
				EndTime:    200 * time.Millisecond,
			}},
			setupMocks: func(mc *minimock.Controller, w *text.WriterMock) {
				w.PrintCandidateMock.Return(nil)
				w.PrintFinalMock.Return(nil)
				w.FinalizeMock.Return(nil)
			},
			expectErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			inputCh := make(chan components.Recognized, 10)
			for _, r := range tt.recognized {
				inputCh <- r
			}
			close(inputCh)
			mc := minimock.NewController(t)
			w := text.NewWriterMock(mc)
			tt.setupMocks(mc, w)

			outputComponent := components.NewOutputComponent(logger.NilLogger, t.TempDir(), io.Discard)

			err := outputComponent.Print(t.Context(), w, inputCh)
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestToFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		outputDir string
		expectErr string
	}{
		{
			name:      "cannot create the file",
			outputDir: string([]byte{0x00}),
			expectErr: "cannot create the file",
		},
		{
			name:      "happy path",
			outputDir: t.TempDir(),
			expectErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			inputCh := make(chan components.Recognized)
			close(inputCh)
			outputComponent := components.NewOutputComponent(logger.NilLogger, tt.outputDir, io.Discard)
			err := outputComponent.ToFile(t.Context(), inputCh)
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestToConsole(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		expectErr string
	}{
		{
			name:      "happy path",
			expectErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			inputCh := make(chan components.Recognized)
			close(inputCh)
			outputComponent := components.NewOutputComponent(logger.NilLogger, "", io.Discard)
			err := outputComponent.ToConsole(t.Context(), inputCh)
			if tt.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
