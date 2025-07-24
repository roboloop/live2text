package storage

import (
	"context"
	"fmt"

	"live2text/internal/services/btt/client/http"
)

const (
	variablePrefix           Key = "LIVE2TEXT_"
	HostVariable                 = variablePrefix + "HOST"
	SelectedDeviceVariable       = variablePrefix + "SELECTED_DEVICE"
	SelectedLanguageVariable     = variablePrefix + "SELECTED_LANGUAGE"
	SelectedFloatingVariable     = variablePrefix + "SELECTED_FLOATING"
	SelectedViewModeVariable     = variablePrefix + "SELECTED_VIEW_MODE"
	TaskIDVariable               = variablePrefix + "TASK_ID"
)

type storage struct {
	httpClient http.Client
}

func NewStorage(httpClient http.Client) Storage {
	return &storage{httpClient: httpClient}
}

func (s *storage) GetValue(ctx context.Context, key Key) (string, error) {
	payload := map[string]string{
		"variableName": string(key),
	}
	text, err := s.httpClient.Send(ctx, "get_string_variable", nil, payload)
	if err != nil {
		return "", fmt.Errorf("cannot get string variable: %w", err)
	}

	return string(text), nil
}

func (s *storage) SetValue(ctx context.Context, key Key, value string) error {
	payload := map[string]string{
		"variableName": string(key),
		"to":           value,
	}

	_, err := s.httpClient.Send(ctx, "set_persistent_string_variable", nil, payload)
	if err != nil {
		return fmt.Errorf("cannot set string variable: %w", err)
	}

	return err
}
