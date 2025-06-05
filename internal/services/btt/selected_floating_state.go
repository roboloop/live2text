package btt

import (
	"context"
	"fmt"
)

func (b *btt) SelectedFloatingState(ctx context.Context) (string, error) {
	value, err := b.getStringVariable(ctx, selectedFloatingStateVariable)
	if err != nil {
		return "", fmt.Errorf("cannot get selected floating state variable: %w", err)
	}

	return value, nil
}
