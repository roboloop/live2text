package btt

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/roboloop/live2text/internal/services/btt/client"
	"github.com/roboloop/live2text/internal/services/btt/client/trigger"
	"github.com/roboloop/live2text/internal/services/btt/storage"
	"github.com/roboloop/live2text/internal/services/btt/tmpl"
	"github.com/roboloop/live2text/internal/services/recognition"
)

var (
	ErrDeviceNotSelected   = errors.New("device is not selected")
	ErrDeviceIsUnavailable = errors.New("device is unavailable")
	ErrLanguageNotSelected = errors.New("language is not selected")
)

type listeningComponent struct {
	logger      *slog.Logger
	recognition recognition.Recognition

	client             client.Client
	storage            storage.Storage
	renderer           tmpl.Renderer
	deviceComponent    DeviceComponent
	languageComponent  LanguageComponent
	viewModeComponent  ViewModeComponent
	floatingComponent  FloatingComponent
	clipboardComponent ClipboardComponent
}

func NewListeningComponent(
	logger *slog.Logger,
	recognition recognition.Recognition,
	client client.Client,
	storage storage.Storage,
	renderer tmpl.Renderer,
	deviceComponent DeviceComponent,
	languageComponent LanguageComponent,
	viewModeComponent ViewModeComponent,
	floatingComponent FloatingComponent,
	clipboardComponent ClipboardComponent,
) ListeningComponent {
	return &listeningComponent{
		logger:             logger,
		recognition:        recognition,
		client:             client,
		storage:            storage,
		renderer:           renderer,
		deviceComponent:    deviceComponent,
		languageComponent:  languageComponent,
		viewModeComponent:  viewModeComponent,
		floatingComponent:  floatingComponent,
		clipboardComponent: clipboardComponent,
	}
}

func (l *listeningComponent) ToggleListening(ctx context.Context) error {
	id, err := l.storage.GetValue(ctx, storage.TaskIDVariable)
	if err != nil {
		return fmt.Errorf("cannot get task id: %w", err)
	}

	if id != "" {
		return l.StopListening(ctx)
	}

	return l.StartListening(ctx)
}

func (l *listeningComponent) StartListening(ctx context.Context) error {
	// Get settings
	device, err := l.deviceComponent.SelectedDevice(ctx)
	if err != nil {
		return fmt.Errorf("cannot get selected device: %w", err)
	}
	if device == "" {
		return ErrDeviceNotSelected
	}
	isAvailable, err := l.deviceComponent.IsAvailable(ctx, device)
	if err != nil {
		return fmt.Errorf("cannot get device availability: %w", err)
	}
	if !isAvailable {
		return ErrDeviceIsUnavailable
	}

	language, err := l.languageComponent.SelectedLanguage(ctx)
	if err != nil {
		return fmt.Errorf("cannot get selected language: %w", err)
	}
	if language == "" {
		return ErrLanguageNotSelected
	}

	// Start recognition
	id, socketPath, err := l.recognition.Start(ctx, device, language)
	if err != nil {
		return fmt.Errorf("cannot start recognition: %w", err)
	}
	if err = l.storage.SetValue(ctx, storage.TaskIDVariable, id); err != nil {
		return fmt.Errorf("cannot set task id: %w", err)
	}

	// Update app scripts that grab subtitles
	t := trigger.NewTrigger().AddCycledScript(l.renderer.ListenSocket(socketPath), defaultInterval)
	if err = l.updateApps(ctx, t); err != nil {
		return err
	}

	// Enable extra features
	if err = l.viewModeComponent.EnableCleanMode(ctx); err != nil {
		return fmt.Errorf("cannot enable clean mode: %w", err)
	}
	if err = l.floatingComponent.ShowFloating(ctx); err != nil {
		return fmt.Errorf("cannot show floating: %w", err)
	}
	if err = l.clipboardComponent.ShowClipboard(ctx); err != nil {
		return fmt.Errorf("cannot show clipboard: %w", err)
	}

	return nil
}

func (l *listeningComponent) StopListening(ctx context.Context) error {
	id, err := l.storage.GetValue(ctx, storage.TaskIDVariable)
	if err != nil {
		return fmt.Errorf("cannot get task id: %w", err)
	}

	if err = l.recognition.Stop(ctx, id); err != nil {
		// do not return until the script has updated
		l.logger.InfoContext(ctx, "Cannot stop task", "error", err)
	}

	if err = l.storage.SetValue(ctx, storage.TaskIDVariable, ""); err != nil {
		return fmt.Errorf("cannot empty task id: %w", err)
	}

	// Update app scripts that grab subtitles
	t := trigger.NewTrigger().AddCycledScript(l.renderer.AppPlaceholder(), 0.0)
	if err = l.updateApps(ctx, t); err != nil {
		return err
	}

	// Disable extra features
	if err = l.viewModeComponent.DisableCleanView(ctx); err != nil {
		return fmt.Errorf("cannot disable clean mode: %w", err)
	}
	if err = l.floatingComponent.HideFloating(ctx); err != nil {
		return fmt.Errorf("cannot hide floating: %w", err)
	}
	if err = l.clipboardComponent.HideClipboard(ctx); err != nil {
		return fmt.Errorf("cannot hide clipboard: %w", err)
	}

	return nil
}

func (l *listeningComponent) updateApps(ctx context.Context, app trigger.Trigger) error {
	if err := l.client.UpdateTrigger(ctx, trigger.TitleApp, app); err != nil {
		return fmt.Errorf("cannot update app: %w", err)
	}
	if err := l.client.UpdateTrigger(ctx, trigger.TitleCleanViewApp, app); err != nil {
		return fmt.Errorf("cannot update clean view app: %w", err)
	}

	return nil
}

func (l *listeningComponent) IsRunning(ctx context.Context) (bool, error) {
	id, err := l.storage.GetValue(ctx, storage.TaskIDVariable)
	if err != nil {
		return false, fmt.Errorf("cannot get task id: %w", err)
	}

	return l.recognition.Has(id), nil
}

func (l *listeningComponent) Text(ctx context.Context) (string, error) {
	id, err := l.storage.GetValue(ctx, storage.TaskIDVariable)
	if err != nil {
		return "", fmt.Errorf("cannot get task id: %w", err)
	}

	text, err := l.recognition.Text(ctx, id)
	if err != nil {
		return "", fmt.Errorf("cannot get text: %w", err)
	}

	return text, nil
}
