package btt

import (
	"context"
	"fmt"

	"github.com/roboloop/live2text/internal/services/btt/client"
	"github.com/roboloop/live2text/internal/services/btt/client/trigger"
	"github.com/roboloop/live2text/internal/services/btt/storage"
)

type settingsComponent struct {
	client  client.Client
	storage storage.Storage
}

func NewSettingsComponent(client client.Client, storage storage.Storage) SettingsComponent {
	return &settingsComponent{client: client, storage: storage}
}

func (s *settingsComponent) SelectSettings(
	ctx context.Context,
	title trigger.Title,
	key storage.Key,
	value string,
) error {
	if err := s.storage.SetValue(ctx, key, value); err != nil {
		return fmt.Errorf("cannot set value: %w", err)
	}

	if err := s.client.RefreshTrigger(ctx, title); err != nil {
		return fmt.Errorf("cannot refresh selected setting: %w", err)
	}

	return nil
}

func (s *settingsComponent) SelectedSetting(ctx context.Context, key storage.Key) (string, error) {
	value, err := s.storage.GetValue(ctx, key)
	if err != nil {
		return "", fmt.Errorf("cannot get selected setting: %w", err)
	}

	return value, nil
}
