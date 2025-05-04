package utils_test

import (
	"context"
	"live2text/internal/utils"
	"testing"
	"time"
)

func TestBroadcaster(t *testing.T) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	ctx = context.Background()

	t.Run("Happy path", func(t *testing.T) {
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()

		value := 42
		inputCh := make(chan int, 1)
		outputs := utils.Broadcaster(ctx, utils.NilLogger, inputCh, 2)
		inputCh <- value
		for i, output := range outputs {
			select {
			case val := <-output:
				if val != value {
					t.Errorf("Output got %v, expected %v", value, val)
				}
			case <-time.After(10 * time.Millisecond):
				t.Errorf("Output [%d]: timed out waiting for value", i)
			}
		}
	})

	t.Run("Discard from channel", func(t *testing.T) {
		t.Skip("broken test")
		ctx, cancel = context.WithCancel(ctx)
		defer cancel()

		inputCh := make(chan int, 1)
		testLogger, testHandler := utils.NewTestLogger()
		outputs := utils.Broadcaster(ctx, testLogger, inputCh, 1)
		output := outputs[0]

		inputCh <- 10
		if val := <-output; val != 10 {
			t.Errorf("Output got %v, expected %v", val, 10)
			return
		}

		inputCh <- 20
		inputCh <- 30
		time.Sleep(100 * time.Millisecond)
		if len(testHandler.Logs) != 1 {
			t.Errorf("Total logs got %v, expected %v", len(testHandler.Logs), 1)
			return
		}
	})
}
