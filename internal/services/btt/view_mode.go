package btt

import (
	"context"
	"fmt"

	"live2text/internal/services/btt/client"
	"live2text/internal/services/btt/client/trigger"
	"live2text/internal/services/btt/storage"
)

type viewModeComponent struct {
	client   client.Client
	settings SettingsComponent
}

func NewViewModeComponent(client client.Client, settings SettingsComponent) ViewModeComponent {
	return &viewModeComponent{client: client, settings: settings}
}

type ViewMode string

const (
	ViewModeClean ViewMode = "Clean"
	ViewModeEmbed ViewMode = "Embed"
)

func (vm ViewMode) isClean() bool {
	return vm == ViewModeClean
}

func (b *viewModeComponent) SelectViewMode(ctx context.Context, viewMode ViewMode) error {
	return b.settings.SelectSettings(
		ctx,
		trigger.TitleSelectedViewMode,
		storage.SelectedViewModeVariable,
		string(viewMode),
	)
}

func (b *viewModeComponent) SelectedViewMode(ctx context.Context) (ViewMode, error) {
	viewMode, err := b.settings.SelectedSetting(ctx, storage.SelectedViewModeVariable)
	return ViewMode(viewMode), err
}

func (b *viewModeComponent) EnableCleanMode(ctx context.Context) error {
	value, err := b.SelectedViewMode(ctx)
	if err != nil {
		return fmt.Errorf("cannot get selected view mode: %w", err)
	}

	if !value.isClean() {
		return nil
	}

	action := trigger.NewOpenDirAction(trigger.TitleCleanViewDir)
	if err = b.client.TriggerAction(ctx, action); err != nil {
		return fmt.Errorf("cannot open dir: %w", err)
	}

	return nil
}

func (b *viewModeComponent) DisableCleanView(ctx context.Context) error {
	action := trigger.NewCloseDirAction()
	if err := b.client.TriggerAction(ctx, action); err != nil {
		return fmt.Errorf("cannot close dir: %w", err)
	}

	return nil
}
