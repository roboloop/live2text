package btt

import (
	"context"
	"fmt"
)

func (b *btt) SelectViewMode(ctx context.Context, viewMode ViewMode) error {
	if err := b.setPersistentStringVariable(ctx, selectedViewModeVariable, string(viewMode)); err != nil {
		return fmt.Errorf("cannot set selected view mode: %w", err)
	}

	selectedViewModeUUID, err := b.getTriggerUUID(ctx, selectedViewModeTitle)
	if err != nil {
		return fmt.Errorf("cannot get selected view mode trigger: %w", err)
	}

	if err = b.RefreshWidget(ctx, selectedViewModeUUID); err != nil {
		return fmt.Errorf("cannot refresh selected view mode widget: %w", err)
	}

	return nil
}
