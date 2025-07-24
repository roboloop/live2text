package btt

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"live2text/internal/services/btt/client"
	"live2text/internal/services/btt/client/trigger"
	"live2text/internal/services/btt/storage"
	"live2text/internal/services/btt/tmpl"
	"live2text/internal/services/recognition"
)

type Floating string

type floatingComponent struct {
	logger      *slog.Logger
	recognition recognition.Recognition

	client   client.Client
	storage  storage.Storage
	renderer tmpl.Renderer
	settings SettingsComponent
}

func NewFloatingComponent(
	logger *slog.Logger,
	recognition recognition.Recognition,
	client client.Client,
	storage storage.Storage,
	renderer tmpl.Renderer,
	settings SettingsComponent,
) FloatingComponent {
	return &floatingComponent{
		logger:      logger,
		recognition: recognition,
		client:      client,
		storage:     storage,
		renderer:    renderer,
		settings:    settings,
	}
}

const (
	FloatingShown  = "Shown"
	FloatingHidden = "Hidden"
)

func (fs Floating) isShown() bool {
	return fs == FloatingShown
}

func (b *floatingComponent) SelectFloating(ctx context.Context, floating Floating) error {
	return b.settings.SelectSettings(
		ctx,
		trigger.TitleSelectedFloating,
		storage.SelectedFloatingVariable,
		string(floating),
	)
}

func (b *floatingComponent) SelectedFloating(ctx context.Context) (Floating, error) {
	floating, err := b.settings.SelectedSetting(ctx, storage.SelectedFloatingVariable)

	return Floating(floating), err
}

func (b *floatingComponent) ShowFloating(ctx context.Context) error {
	value, err := b.SelectedFloating(ctx)
	if err != nil {
		return fmt.Errorf("cannot get selected floating: %w", err)
	}

	if !value.isShown() {
		return nil
	}

	action := trigger.NewOpenFloatingAction(trigger.TitleStreamingText)
	if err = b.client.TriggerAction(ctx, action); err != nil {
		return fmt.Errorf("cannot show floating: %w", err)
	}

	return nil
}

func (b *floatingComponent) HideFloating(ctx context.Context) error {
	action := trigger.NewCloseFloatingAction(trigger.TitleStreamingText)
	if err := b.client.TriggerAction(ctx, action); err != nil {
		return fmt.Errorf("cannot hide floating: %w", err)
	}

	return nil
}

func (b *floatingComponent) FloatingPage() string {
	return b.renderer.FloatingPage()
}

func (b *floatingComponent) StreamText(ctx context.Context) (<-chan string, <-chan error, error) {
	id, err := b.storage.GetValue(ctx, storage.TaskIDVariable)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot get listening to socket variable: %w", err)
	}

	t := time.NewTicker(time.Duration(defaultInterval * float64(time.Second)))
	textCh := make(chan string, 1024)
	errCh := make(chan error, 1)

	go func() {
		defer t.Stop()
		defer close(textCh)
		defer close(errCh)

		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				text, errText := b.recognition.Text(ctx, id)
				if errText != nil {
					b.logger.ErrorContext(ctx, "Cannot get text", "error", errText)
					errCh <- errText
					return
				}
				textCh <- text
			}
		}
	}()

	return textCh, errCh, nil
}
