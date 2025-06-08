package btt

import (
	"context"
	"fmt"
)

func (b *btt) SelectedViewMode(ctx context.Context) (string, error) {
	value, err := b.getStringVariable(ctx, selectedViewModeVariable)
	if err != nil {
		return "", fmt.Errorf("cannot get selected view mode variable: %w", err)
	}

	return value, nil
}
