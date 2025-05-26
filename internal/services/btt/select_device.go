package btt

import (
	"context"
	"fmt"
)

func (b *btt) SelectDevice(ctx context.Context, device string) error {
	if err := b.setPersistentStringVariable(ctx, selectedDeviceVariable, device); err != nil {
		return fmt.Errorf("cannot set selected device: %w", err)
	}

	selectedDeviceUUID, err := b.getTriggerUUID(ctx, selectedDeviceTitle)
	if err != nil {
		return fmt.Errorf("cannot get selected device trigger: %w", err)
	}

	if err = b.RefreshWidget(ctx, selectedDeviceUUID); err != nil {
		return fmt.Errorf("cannot refresh selected device widget: %w", err)
	}

	return nil
}
