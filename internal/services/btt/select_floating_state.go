package btt

import (
	"context"
	"fmt"
)

func (b *btt) SelectFloatingState(ctx context.Context, floatingState FloatingState) error {
	if err := b.setPersistentStringVariable(ctx, selectedFloatingStateVariable, string(floatingState)); err != nil {
		return fmt.Errorf("cannot set selected floating state: %w", err)
	}

	selectedFloatingStateUUID, err := b.getTriggerUUID(ctx, selectedFloatingStateTitle)
	if err != nil {
		return fmt.Errorf("cannot get selected floating state trigger: %w", err)
	}

	if err = b.RefreshWidget(ctx, selectedFloatingStateUUID); err != nil {
		return fmt.Errorf("cannot refresh selected floating state widget: %w", err)
	}

	return nil
}
