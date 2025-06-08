package btt

import (
	"context"
	"fmt"
)

func (b *btt) IsRunning(ctx context.Context) (bool, error) {
	id, err := b.getStringVariable(ctx, taskIDVariable)
	if err != nil {
		return false, fmt.Errorf("cannot get listeting socket variable: %w", err)
	}

	return b.recognition.Has(id), nil
}
