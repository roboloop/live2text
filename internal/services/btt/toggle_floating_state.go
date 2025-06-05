package btt

import (
	"context"
	"fmt"

	"live2text/internal/services/btt/payload"
)

type FloatingState string

const (
	FloatingStateShown  = "Shown"
	FloatingStateHidden = "Hidden"
)

func (fs FloatingState) isShown() bool {
	return fs == FloatingStateShown
}

func (b *btt) showFloatingState(ctx context.Context) error {
	value, err := b.getStringVariable(ctx, selectedFloatingStateVariable)
	if err != nil {
		return fmt.Errorf("cannot get selected floating state variable: %w", err)
	}

	if !FloatingState(value).isShown() {
		return nil
	}

	actionPayload := make(payload.Payload).AddMap(map[string]any{
		"BTTPredefinedActionType": payload.ActionTypeOpenFloatingMenu,
		"BTTAdditionalActionData": map[string]any{
			"BTTMenuActionMenuID": floatingStateGroupTitle,
		},
	})

	if _, err = b.httpClient.Send(ctx, "trigger_action", actionPayload, nil); err != nil {
		return fmt.Errorf("cannot trigger show floating state trigger: %w", err)
	}

	return nil
}

func (b *btt) hideFloatingState(ctx context.Context) error {
	actionPayload := make(payload.Payload).AddMap(map[string]any{
		"BTTPredefinedActionType": payload.ActionTypeCloseFloatingMenu,
		"BTTAdditionalActionData": map[string]any{
			"BTTMenuActionMenuID": floatingStateGroupTitle,
		},
	})

	if _, err := b.httpClient.Send(ctx, "trigger_action", actionPayload, nil); err != nil {
		return fmt.Errorf("cannot trigger hide floating state trigger: %w", err)
	}

	return nil
}
