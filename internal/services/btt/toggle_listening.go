package btt

import (
	"context"
	"fmt"

	"live2text/internal/services/btt/payload"
)

func (b *btt) ToggleListening(ctx context.Context) error {
	id, err := b.getStringVariable(ctx, taskIDVariable)
	if err != nil {
		return fmt.Errorf("cannot get listeting socket variable: %w", err)
	}

	if id != "" {
		return b.StopListening(ctx, id)
	}

	return b.StartListening(ctx)
}

func (b *btt) StopListening(ctx context.Context, id string) error {
	var err error
	if err = b.setPersistentStringVariable(ctx, taskIDVariable, ""); err != nil {
		return fmt.Errorf("cannot remove task id: %w", err)
	}

	if err = b.recognition.Stop(ctx, id); err != nil {
		b.logger.InfoContext(ctx, "Cannot stop task", "error", err)
	}

	rendered, err := b.renderer.Render("listen_socket", map[string]any{})
	if err != nil {
		return fmt.Errorf("cannot render listen_socket script: %w", err)
	}

	appPayload := make(payload.Payload).
		AddShell(rendered, 0.0, payload.ShellTypeNone).
		AddMap(map[string]any{
			"BTTTriggerConfig": map[string]any{
				"BTTTouchBarScriptUpdateInterval": 0.0,
			},
		})

	appUUID, cleanViewAppUUID, err := b.appTriggers(ctx)
	if err != nil {
		return fmt.Errorf("cannot get app triggers: %w", err)
	}
	if _, err = b.httpClient.Send(ctx, "update_trigger", appPayload, map[string]string{"uuid": appUUID}); err != nil {
		return fmt.Errorf("cannot update app trigger: %w", err)
	}
	if _, err = b.httpClient.Send(ctx, "update_trigger", appPayload, map[string]string{"uuid": cleanViewAppUUID}); err != nil {
		return fmt.Errorf("cannot update clean view app trigger: %w", err)
	}

	if err = b.disableCleanView(ctx); err != nil {
		return fmt.Errorf("cannot disable clean view: %w", err)
	}

	if err = b.hideFloatingState(ctx); err != nil {
		return fmt.Errorf("cannot hide floating state: %w", err)
	}

	return nil
}

func (b *btt) StartListening(ctx context.Context) error {
	var (
		device     string
		language   string
		id         string
		socketPath string
		err        error
	)
	if device, err = b.getStringVariable(ctx, selectedDeviceVariable); err != nil {
		return fmt.Errorf("cannot get selected device: %w", err)
	}
	if language, err = b.getStringVariable(ctx, selectedLanguageVariable); err != nil {
		return fmt.Errorf("cannot get selected language: %w", err)
	}

	if id, socketPath, err = b.recognition.Start(ctx, device, language); err != nil {
		return fmt.Errorf("cannot start recognition: %w", err)
	}
	if err = b.setPersistentStringVariable(ctx, taskIDVariable, id); err != nil {
		return fmt.Errorf("cannot set task id variable: %w", err)
	}

	rendered, err := b.renderer.Render("listen_socket", map[string]any{"SocketPath": socketPath})
	if err != nil {
		return fmt.Errorf("cannot render print_status script: %w", err)
	}

	appPayload := make(payload.Payload).
		AddShell(rendered, b.interval, payload.ShellTypeNone).
		AddMap(map[string]any{
			"BTTTriggerConfig": map[string]any{
				"BTTTouchBarScriptUpdateInterval": defaultInterval,
			},
		})

	appUUID, cleanViewAppUUID, err := b.appTriggers(ctx)
	if err != nil {
		return fmt.Errorf("cannot get app triggers: %w", err)
	}
	if _, err = b.httpClient.Send(ctx, "update_trigger", appPayload, map[string]string{"uuid": appUUID}); err != nil {
		return fmt.Errorf("cannot update app trigger: %w", err)
	}
	if _, err = b.httpClient.Send(ctx, "update_trigger", appPayload, map[string]string{"uuid": cleanViewAppUUID}); err != nil {
		return fmt.Errorf("cannot update clean view app trigger: %w", err)
	}

	if err = b.enableCleanMode(ctx); err != nil {
		return fmt.Errorf("cannot enable clean mode: %w", err)
	}

	if err = b.showFloatingState(ctx); err != nil {
		return fmt.Errorf("cannot show floating state: %w", err)
	}

	return nil
}

func (b *btt) appTriggers(ctx context.Context) (string, string, error) {
	appTrigger, err := b.getTrigger(ctx, appTitle)
	if err != nil {
		return "", "", fmt.Errorf("cannot get app trigger: %w", err)
	}
	cleanViewAppTrigger, err := b.getTrigger(ctx, cleanViewAppTitle)
	if err != nil {
		return "", "", fmt.Errorf("cannot get clean view app trigger: %w", err)
	}

	appUUID := appTrigger[bttUUID].(string)
	cleanViewAppUUID := cleanViewAppTrigger[bttUUID].(string)

	return appUUID, cleanViewAppUUID, nil
}
