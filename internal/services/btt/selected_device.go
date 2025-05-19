package btt

import (
	"context"
	"fmt"
)

func (b *btt) SelectedDevice(ctx context.Context) (string, error) {
	value, err := b.getStringVariable(ctx, selectedDeviceVariable)
	if err != nil {
		return "", fmt.Errorf("cannot get selected device variable: %w", err)
	}

	return value, nil
}
