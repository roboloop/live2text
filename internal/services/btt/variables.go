package btt

import (
	"context"
	"fmt"
)

const (
	variablePrefix                = "LIVE2TEXT_"
	hostVariable                  = variablePrefix + "HOST"
	selectedDeviceVariable        = variablePrefix + "SELECTED_DEVICE"
	selectedLanguageVariable      = variablePrefix + "SELECTED_LANGUAGE"
	selectedFloatingStateVariable = variablePrefix + "SELECTED_FLOATING_STATE"
	selectedViewModeVariable      = variablePrefix + "SELECTED_VIEW_MODE"
	taskIDVariable                = variablePrefix + "TASK_ID"
)

func (b *btt) getStringVariable(ctx context.Context, variable string) (string, error) {
	payload := map[string]string{
		"variableName": variable,
	}
	text, err := b.httpClient.Send(ctx, "get_string_variable", nil, payload)
	if err != nil {
		return "", fmt.Errorf("cannot get string variable: %w", err)
	}

	return string(text), nil
}

func (b *btt) setPersistentStringVariable(ctx context.Context, variable string, to string) error {
	payload := map[string]string{
		"variableName": variable,
		"to":           to,
	}

	_, err := b.httpClient.Send(ctx, "set_persistent_string_variable", nil, payload)
	if err != nil {
		return fmt.Errorf("cannot set string variable: %w", err)
	}

	return err
}
