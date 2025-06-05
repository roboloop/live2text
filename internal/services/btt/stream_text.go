package btt

import (
	"context"
	"fmt"
	"time"
)

func (b *btt) StreamText(ctx context.Context) (<-chan string, <-chan error, error) {
	id, err := b.getStringVariable(ctx, taskIDVariable)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot get listeting socket variable: %w", err)
	}

	t := time.NewTicker(250 * time.Millisecond)
	textCh := make(chan string, 1024)
	errCh := make(chan error, 1)

	go func() {
		defer t.Stop()
		defer close(textCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				text, errSubs := b.recognition.Subs(ctx, id)
				if errSubs != nil {
					b.logger.ErrorContext(ctx, "Cannot get subs", "error", errSubs)
					errCh <- errSubs
					return
				}
				textCh <- text
			}
		}
	}()

	return textCh, errCh, nil
}
