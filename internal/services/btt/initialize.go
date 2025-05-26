package btt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"live2text/internal/services/btt/payload"
	"net/url"
	"strings"
)

func (b *btt) Initialize(ctx context.Context) error {
	if err := b.setPersistentStringVariable(ctx, hostVariable, b.appAddress); err != nil {
		return fmt.Errorf("cannot set app address: %w", err)
	}

	if err := b.addSettingsMenu(ctx); err != nil {
		return fmt.Errorf("cannot settings menu: %w", err)
	}

	if err := b.addMainMenu(ctx); err != nil {
		return fmt.Errorf("cannot add main menu: %w", err)
	}

	return nil
}

func (b *btt) addSettingsMenu(ctx context.Context) error {
	settingsPayload := make(
		payload.Payload,
	).AddTrigger(settingsTitle, payload.TriggerDirectory, payload.TouchBar, payload.ActionTypeExecuteScript, true)

	settingsUUID, err := b.addTrigger(ctx, settingsPayload, payload.MainOrderSettings, "")
	if err != nil {
		return fmt.Errorf("cannot create settings group: %w", err)
	}

	closePayload := make(payload.Payload).
		AddTrigger("Close Group", payload.TriggerTouchBarButton, payload.TouchBar, payload.ActionTypeCloseGroup, false).
		AddIcon("xmark.circle.fill", 25, true)
	if _, err = b.addTrigger(ctx, closePayload, payload.SettingsOrderCloseGroup, settingsUUID); err != nil {
		return fmt.Errorf("cannot create close trigger: %w", err)
	}

	rendered, err := b.renderer.Render("print_status", map[string]any{"AppAddress": b.appAddress})
	if err != nil {
		return fmt.Errorf("cannot render print_status script: %w", err)
	}
	statusPayload := make(payload.Payload).
		AddTrigger("⏳", payload.TriggerShellScript, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
		AddShell(rendered, 15, payload.ShellTypeNone)
	if _, err = b.addTrigger(ctx, statusPayload, payload.SettingsOrderStatus, settingsUUID); err != nil {
		return fmt.Errorf("cannot create status trigger: %w", err)
	}

	if err = b.addDeviceMenu(ctx, settingsUUID); err != nil {
		return fmt.Errorf("cannot create device menu: %w", err)
	}

	if err = b.addLanguageMenu(ctx, settingsUUID); err != nil {
		return fmt.Errorf("cannot create language menu: %w", err)
	}

	return nil
}

func (b *btt) addDeviceMenu(ctx context.Context, parentUUID string) error {
	devicePayload := make(payload.Payload).
		AddTrigger(deviceGroupTitle, payload.TriggerDirectory, payload.TouchBar, payload.ActionTypeExecuteScript, false).
		AddIcon("microphone", 22, false)
	deviceUUID, err := b.addTrigger(ctx, devicePayload, payload.SettingsOrderDevice, parentUUID)
	if err != nil {
		return fmt.Errorf("cannot create device group: %w", err)
	}

	closePayload := make(payload.Payload).
		AddClose(settingsTitle).
		AddIcon("xmark.circle.fill", 25, true)
	if _, err = b.addTrigger(ctx, closePayload, payload.DeviceOrderCloseGroup, deviceUUID); err != nil {
		return fmt.Errorf("cannot create close trigger: %w", err)
	}

	rendered, err := b.renderer.Render("print_selected_device", map[string]any{"AppAddress": b.appAddress})
	if err != nil {
		return fmt.Errorf("cannot render print_status script: %w", err)
	}
	selectedDevicePayload := make(payload.Payload).
		AddTrigger(selectedDeviceTitle, payload.TriggerShellScript, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
		AddShell(rendered, 15, payload.ShellTypeNone).
		AddMap(map[string]any{"BTTTriggerConfig": map[string]any{"BTTTouchBarFreeSpaceAfterButton": 25}})
	if _, err = b.addTrigger(ctx, selectedDevicePayload, payload.DeviceOrderSelectedDevice, deviceUUID); err != nil {
		return fmt.Errorf("cannot create close trigger: %w", err)
	}

	return nil
}

func (b *btt) addLanguageMenu(ctx context.Context, parentUUID string) error {
	devicePayload := make(payload.Payload).
		AddTrigger(languageGroupTitle, payload.TriggerDirectory, payload.TouchBar, payload.ActionTypeExecuteScript, false).
		AddIcon("character", 22, false)
	languageUUID, err := b.addTrigger(ctx, devicePayload, payload.SettingsOrderLanguage, parentUUID)
	if err != nil {
		return fmt.Errorf("cannot create language group: %w", err)
	}

	closePayload := make(payload.Payload).
		AddClose(settingsTitle).
		AddIcon("xmark.circle.fill", 25, true)
	if _, err = b.addTrigger(ctx, closePayload, payload.LanguageOrderCloseGroup, languageUUID); err != nil {
		return fmt.Errorf("cannot create close trigger: %w", err)
	}

	rendered, err := b.renderer.Render("print_selected_language", map[string]any{"AppAddress": b.appAddress})
	if err != nil {
		return fmt.Errorf("cannot render print_selected_language script: %w", err)
	}
	selectedDevicePayload := make(payload.Payload).
		AddTrigger(selectedLanguageTitle, payload.TriggerShellScript, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
		AddShell(rendered, 15, payload.ShellTypeNone).
		AddMap(map[string]any{"BTTTriggerConfig": map[string]any{"BTTTouchBarFreeSpaceAfterButton": 25}})
	if _, err = b.addTrigger(ctx, selectedDevicePayload, payload.LanguageOrderSelectedLanguage, languageUUID); err != nil {
		return fmt.Errorf("cannot create close trigger: %w", err)
	}

	for i, language := range b.languages {
		rendered, err = b.renderer.Render(
			"select_language",
			map[string]any{"AppAddress": b.appAddress, "Language": language},
		)
		if err != nil {
			return fmt.Errorf("cannot render print_selected_language script: %w", err)
		}

		languagePayload := make(payload.Payload).
			AddTrigger(language, payload.TriggerTouchBarButton, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
			AddShell(rendered, 0, payload.ShellTypeAdditional)
		if _, err = b.addTrigger(ctx, languagePayload, payload.LanguageOrderSelectedLanguage+payload.Order(1+i), languageUUID); err != nil {
			return fmt.Errorf("cannot create close trigger: %w", err)
		}
	}

	return nil
}

func (b *btt) addMainMenu(ctx context.Context) error {
	// TODO: code is repeated
	jsonPayload := map[string]any{
		"BTTPredefinedActionType": payload.ActionTypeOpenGroup,
		"BTTOpenGroupWithName":    settingsTitle,
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(jsonPayload); err != nil {
		return fmt.Errorf("cannot ecnode json payload: %w", err)
	}
	query := url.Values{}
	query.Set("json", buf.String())
	encoded := strings.ReplaceAll(query.Encode(), "+", "%20")

	rendered, err := b.renderer.Render(
		"open_settings",
		map[string]any{"AppAddress": b.appAddress, "BttAddress": b.bttAddress, "Query": encoded},
	)
	if err != nil {
		return fmt.Errorf("cannot render print_selected_language script: %w", err)
	}

	namedPayload := make(payload.Payload).
		AddTrigger(settingsTitle, payload.TriggerNamed, payload.OtherTriggers, payload.ActionTypeEmptyPlaceholder, false).
		AddShell(rendered, 0, payload.ShellTypeAdditional)
	if _, err = b.addTrigger(ctx, namedPayload, payload.NamedOrderSettings, ""); err != nil {
		return fmt.Errorf("cannot create named trigger: %w", err)
	}

	toggleRendered, err := b.renderer.Render("toggle", map[string]any{"AppAddress": b.appAddress})
	if err != nil {
		return fmt.Errorf("cannot render toggle script: %w", err)
	}
	listenSocketRendered, err := b.renderer.Render("listen_socket", map[string]any{})
	if err != nil {
		return fmt.Errorf("cannot render listen_socket script: %w", err)
	}

	mainPayload := make(payload.Payload).
		AddTrigger(appTitle, payload.TriggerShellScript, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
		AddShell(toggleRendered, 0, payload.ShellTypeEmbed).
		AddShell(listenSocketRendered, 0, payload.ShellTypeNone).
		AddMap(map[string]any{
			"BTTTriggerConfig": map[string]any{
				"BTTTouchBarLongPressActionName": settingsTitle,
				"BTTTouchBarButtonFontSize":      12,
				"BTTTouchBarButtonColor":         "0.0, 0.0, 0.0, 255.0",
				"BTTTouchBarButtonTextAlignment": 0,
			},
		})
	if _, err = b.addTrigger(ctx, mainPayload, payload.MainOrderSubs, ""); err != nil {
		return fmt.Errorf("cannot create main trigger: %w", err)
	}

	return nil
}
