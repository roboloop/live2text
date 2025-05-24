package utils

import (
	"context"
	"log/slog"
)

func Broadcaster[T any](ctx context.Context, logger *slog.Logger, input <-chan T, total int) []<-chan T {
	logger = logger.With("service", "Broadcaster")
	var (
		outputs []chan T
		results []<-chan T
	)
	for i := 0; i < total; i++ {
		ch := make(chan T, cap(input))
		outputs = append(outputs, ch)
		results = append(results, ch)
	}

	go func() {
		defer func() {
			logger.ErrorContext(ctx, "Shutting down...")
			for _, output := range outputs {
				close(output)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-input:
				for i, output := range outputs {
					select {
					case output <- msg:
						continue
					default:
						logger.ErrorContext(ctx, "Message dropped", "output", i)
					}
				}
			}
		}
	}()

	return results
}
