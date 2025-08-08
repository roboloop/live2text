package utils_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/utils"
	"github.com/roboloop/live2text/internal/utils/logger"
)

func TestBroadcaster(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		inputCh := make(chan int, 1)
		outputs := utils.Broadcaster(ctx, logger.NilLogger, inputCh, []string{"test1", "test2"})

		value := 42
		inputCh <- value
		for _, output := range outputs {
			select {
			case v := <-output:
				require.Equal(t, value, v)
			case <-time.After(10 * time.Millisecond):
				t.Fatal("timed out waiting for value")
			}
		}
	})

	t.Run("drops when blocked", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		inputCh := make(chan int, 1)
		testLogger, testHandler := logger.NewCaptureLogger()
		outputs := utils.Broadcaster(ctx, testLogger, inputCh, []string{"test1"})
		output := outputs[0]

		inputCh <- 10
		require.Equal(t, 10, <-output)

		inputCh <- 20
		inputCh <- 30
		time.Sleep(20 * time.Millisecond)

		require.Len(t, testHandler.Logs, 1)
		require.Contains(t, testHandler.Logs[0].Msg, "Message dropped")

		require.Len(t, testHandler.Logs[0].Attrs, 2)
		require.Contains(t, testHandler.Logs[0].Attrs[1].String(), "name=test1")
	})
}
