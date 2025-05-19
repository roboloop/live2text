package btt

import (
	"context"
	"fmt"
	"live2text/internal/services/btt/payload"
)

func (b *btt) ToggleListening(ctx context.Context) error {
	id, err := b.getStringVariable(ctx, taskIdVariable)
	if err != nil {
		return fmt.Errorf("cannot get listeting socket variable: %w", err)
	}

	appTrigger, err := b.getTrigger(ctx, appTitle)
	if err != nil {
		return fmt.Errorf("cannot get app trigger: %w", err)
	}
	uuid := appTrigger[bttUuid].(string)

	if id != "" {
		return b.stop(ctx, uuid, id)
	}

	return b.start(ctx, uuid)
}

func (b *btt) stop(ctx context.Context, uuid, id string) error {
	var err error
	if err = b.setPersistentStringVariable(ctx, taskIdVariable, ""); err != nil {
		return fmt.Errorf("cannot remove task id")
	}

	if err = b.recognition.Stop(ctx, id); err != nil {
		b.logger.ErrorContext(ctx, "Cannot stop task", "error", err)
	}

	rendered, err := b.renderer.Render("listen_socket", map[string]any{})
	if err != nil {
		return fmt.Errorf("cannot render listen_socket script: %w", err)
	}

	appPayload := make(payload.Payload).
		AddShell(rendered, 0.0, payload.ShellTypeNone).
		AddMap(map[string]any{
			"BTTTriggerConfig": map[string]any{
				//"BTTTouchBarAppleScriptStringRunOnInit": 0,
				"BTTTouchBarScriptUpdateInterval": 0.0,
			},
		})

	if _, err = b.httpClient.Send(ctx, "update_trigger", appPayload, map[string]string{"uuid": uuid}); err != nil {
		return fmt.Errorf("cannot update app trigger: %w", err)
	}

	return nil
}

func (b *btt) start(ctx context.Context, uuid string) error {
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
	if err = b.setPersistentStringVariable(ctx, taskIdVariable, id); err != nil {
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
				//"BTTTouchBarAppleScriptStringRunOnInit": 1,
				"BTTTouchBarScriptUpdateInterval": defaultInterval,
			},
		})

	if _, err = b.httpClient.Send(ctx, "update_trigger", appPayload, map[string]string{"uuid": uuid}); err != nil {
		return fmt.Errorf("cannot update app trigger: %w", err)
	}

	return nil
}
