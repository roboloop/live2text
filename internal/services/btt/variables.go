package btt

import (
	"context"
	"fmt"
)

const (
	hostVariable     = "LIVE2TEXT_HOST"
	scriptVariable   = "LIVE2TEXT_SCRIPT"
	subsUUIDVariable = "LIVE2TEXT_SUBS_UUID"

	deviceUUIDVariable         = "LIVE2TEXT_DEVICE_UUID"
	selectedDeviceVariable     = "LIVE2TEXT_SELECTED_DEVICE"
	selectedDeviceUUIDVariable = "LIVE2TEXT_SELECTED_DEVICE_UUID"

	languageUUIDVariable         = "LIVE2TEXT_LANGUAGE_UUID"
	selectedLanguageVariable     = "LIVE2TEXT_SELECTED_LANGUAGE"
	selectedLanguageUUIDVariable = "LIVE2TEXT_SELECTED_LANGUAGE_UUID"

	listeningSocketVariable = "LIVE2TEXT_LISTENING_SOCKET"
	taskIdVariable          = "LIVE2TEXT_TASK_ID"
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
