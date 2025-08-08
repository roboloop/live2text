package btt

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/roboloop/live2text/internal/services/btt/client"
	"github.com/roboloop/live2text/internal/services/btt/client/trigger"
	"github.com/roboloop/live2text/internal/services/btt/storage"
)

type Clipboard string

func (c Clipboard) isShown() bool {
	return c == ClipboardShown
}

const (
	ClipboardShown  Clipboard = "Shown"
	ClipboardHidden Clipboard = "Hidden"
)

type clipboardComponent struct {
	logger   *slog.Logger
	client   client.Client
	settings SettingsComponent
}

func NewClipboardComponent(logger *slog.Logger, client client.Client, settings SettingsComponent) ClipboardComponent {
	return &clipboardComponent{logger: logger, client: client, settings: settings}
}

func (c *clipboardComponent) SelectClipboard(ctx context.Context, clipboard Clipboard) error {
	return c.settings.SelectSettings(
		ctx,
		trigger.TitleSelectedClipboard,
		storage.SelectedClipboardVariable,
		string(clipboard),
	)
}

func (c *clipboardComponent) SelectedClipboard(ctx context.Context) (Clipboard, error) {
	clipboard, err := c.settings.SelectedSetting(ctx, storage.SelectedClipboardVariable)

	return Clipboard(clipboard), err
}

func (c *clipboardComponent) ShowClipboard(ctx context.Context) error {
	value, err := c.SelectedClipboard(ctx)
	if err != nil {
		return fmt.Errorf("cannot get selected clipboard: %w", err)
	}

	if !value.isShown() {
		return nil
	}

	if err = c.updateClipboards(ctx, trigger.NewTrigger().AddEnabled()); err != nil {
		return fmt.Errorf("cannot update clipboards: %w", err)
	}

	return nil
}

func (c *clipboardComponent) HideClipboard(ctx context.Context) error {
	if err := c.updateClipboards(ctx, trigger.NewTrigger().AddDisabled()); err != nil {
		return fmt.Errorf("cannot update clipboard: %w", err)
	}

	return nil
}

func (c *clipboardComponent) updateClipboards(ctx context.Context, patch trigger.Trigger) error {
	err1 := c.client.UpdateTrigger(ctx, trigger.TitleClipboard, patch)
	err2 := c.client.UpdateTrigger(ctx, trigger.TitleCleanViewClipboard, patch)

	return errors.Join(err1, err2)
}
