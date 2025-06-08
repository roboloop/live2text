package btt

import (
	"context"
	"fmt"

	"live2text/internal/services/btt/payload"
)

type ViewMode string

const (
	ViewModeClean = "Clean"
	ViewModeEmbed = "Embed"
)

func (vm ViewMode) isClean() bool {
	return vm == ViewModeClean
}

func (b *btt) enableCleanMode(ctx context.Context) error {
	value, err := b.getStringVariable(ctx, selectedViewModeVariable)
	if err != nil {
		return fmt.Errorf("cannot get selected view mode variable: %w", err)
	}

	if !ViewMode(value).isClean() {
		return nil
	}

	actionPayload := make(payload.Payload).AddMap(map[string]any{
		"BTTPredefinedActionType": payload.ActionTypeOpenGroup,
		"BTTOpenGroupWithName":    cleanViewTitle,
	})

	if _, err = b.httpClient.Send(ctx, "trigger_action", actionPayload, nil); err != nil {
		return fmt.Errorf("cannot trigger open group: %w", err)
	}

	return nil
}

func (b *btt) disableCleanView(ctx context.Context) error {
	actionPayload := make(payload.Payload).AddMap(map[string]any{
		"BTTPredefinedActionType": payload.ActionTypeCloseGroup,
	})

	if _, err := b.httpClient.Send(ctx, "trigger_action", actionPayload, nil); err != nil {
		return fmt.Errorf("cannot trigger close group: %w", err)
	}

	return nil
}
