package utils

import (
	"context"
	"log/slog"
)

func Broadcaster[T any](ctx context.Context, logger *slog.Logger, input <-chan T, names []string) []<-chan T {
	logger = logger.With("component", "broadcaster")
	var (
		outputs []chan T
		results []<-chan T
	)
	for range names {
		ch := make(chan T, cap(input))
		outputs = append(outputs, ch)
		results = append(results, ch)
	}

	go func() {
		defer func() {
			logger.InfoContext(ctx, "Shutting down...")
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
						logger.ErrorContext(ctx, "Message dropped", "name", names[i])
					}
				}
			}
		}
	}()

	return results
}
