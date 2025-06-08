package btt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"live2text/internal/services/btt/payload"
	"live2text/internal/services/metrics"
)

func (b *btt) Initialize(ctx context.Context) error {
	if err := b.setPersistentStringVariable(ctx, hostVariable, b.appAddress); err != nil {
		return fmt.Errorf("cannot set app address: %w", err)
	}

	if err := b.addSettingsSection(ctx); err != nil {
		return fmt.Errorf("cannot add settings section: %w", err)
	}

	if err := b.addCleanModeViewSection(ctx); err != nil {
		return fmt.Errorf("cannot add clean mode view section: %w", err)
	}

	if err := b.addMainSection(ctx); err != nil {
		return fmt.Errorf("cannot add main section: %w", err)
	}

	return nil
}

func (b *btt) addSettingsSection(ctx context.Context) error {
	settingsPayload := make(
		payload.Payload,
	).AddTrigger(settingsTitle, payload.TriggerDirectory, payload.TouchBar, payload.ActionTypeNone, true)

	settingsUUID, err := b.addTrigger(ctx, settingsPayload, payload.MainOrderSettings, "")
	if err != nil {
		return fmt.Errorf("cannot create settings group: %w", err)
	}

	var openClenViewEncoded string
	var closeGroupEncoded string
	if openClenViewEncoded, err = encodeForScript(map[string]any{
		"BTTPredefinedActionType": payload.ActionTypeOpenGroup,
		"BTTOpenGroupWithName":    cleanViewTitle,
	}); err != nil {
		return err
	}
	if closeGroupEncoded, err = encodeForScript(map[string]any{"BTTPredefinedActionType": payload.ActionTypeCloseGroup}); err != nil {
		return err
	}

	rendered, err := b.renderer.Render("close_settings", map[string]any{
		"AppAddress":         b.appAddress,
		"BttAddress":         b.bttAddress,
		"CloseGroupQuery":    closeGroupEncoded,
		"OpenCleanViewQuery": openClenViewEncoded,
		"CleanViewMode":      ViewModeClean,
	})
	if err != nil {
		return fmt.Errorf("cannot render close_settings script: %w", err)
	}

	closeSettingsPayload := make(payload.Payload).
		AddTrigger("Close Group", payload.TriggerTouchBarButton, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
		AddShell(rendered, payload.IntervalNone, payload.ShellTypeEmbed).
		AddIcon("xmark.circle.fill", 25, true)
	if _, err = b.addTrigger(ctx, closeSettingsPayload, payload.SettingsOrderCloseGroup, settingsUUID); err != nil {
		return fmt.Errorf("cannot create close settings trigger: %w", err)
	}

	rendered, err = b.renderer.Render("print_status", map[string]any{"AppAddress": b.appAddress})
	if err != nil {
		return fmt.Errorf("cannot render print_status script: %w", err)
	}
	statusPayload := make(payload.Payload).
		AddTrigger("⏳", payload.TriggerShellScript, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
		AddShell(rendered, payload.IntervalDefault, payload.ShellTypeNone)
	if _, err = b.addTrigger(ctx, statusPayload, payload.SettingsOrderStatus, settingsUUID); err != nil {
		return fmt.Errorf("cannot create status trigger: %w", err)
	}

	if err = b.addDeviceSection(ctx, settingsUUID); err != nil {
		return fmt.Errorf("cannot create device section: %w", err)
	}

	if err = b.addLanguageSection(ctx, settingsUUID); err != nil {
		return fmt.Errorf("cannot create language section: %w", err)
	}

	if err = b.addViewModeSection(ctx, settingsUUID); err != nil {
		return fmt.Errorf("cannot create view mode section: %w", err)
	}

	if err = b.addFloatingSection(ctx, settingsUUID); err != nil {
		return fmt.Errorf("cannot create floating section: %w", err)
	}

	if err = b.addMetricsSection(ctx, settingsUUID); err != nil {
		return fmt.Errorf("cannot create metrics section: %w", err)
	}

	return nil
}

func (b *btt) addDeviceSection(ctx context.Context, parentUUID string) error {
	devicePayload := make(payload.Payload).
		AddTrigger(deviceGroupTitle, payload.TriggerDirectory, payload.TouchBar, payload.ActionTypeExecuteScript, false).
		AddIcon(payload.IconMicrophone, payload.DefaultHeight, false)
	deviceUUID, err := b.addTrigger(ctx, devicePayload, payload.SettingsOrderDevice, parentUUID)
	if err != nil {
		return fmt.Errorf("cannot create device group: %w", err)
	}

	closePayload := make(payload.Payload).AddClose(settingsTitle)
	if _, err = b.addTrigger(ctx, closePayload, payload.DeviceOrderCloseGroup, deviceUUID); err != nil {
		return fmt.Errorf("cannot create selected device trigger: %w", err)
	}

	rendered, err := b.renderer.Render("print_selected_device", map[string]any{"AppAddress": b.appAddress})
	if err != nil {
		return fmt.Errorf("cannot render print_status script: %w", err)
	}
	selectedDevicePayload := make(payload.Payload).
		AddTrigger(selectedDeviceTitle, payload.TriggerShellScript, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
		AddShell(rendered, payload.IntervalDefault, payload.ShellTypeNone).
		AddMap(map[string]any{"BTTTriggerConfig": map[string]any{"BTTTouchBarFreeSpaceAfterButton": 25}})
	if _, err = b.addTrigger(ctx, selectedDevicePayload, payload.DeviceOrderSelectedDevice, deviceUUID); err != nil {
		return fmt.Errorf("cannot create select device trigger: %w", err)
	}

	return nil
}

func (b *btt) addLanguageSection(ctx context.Context, parentUUID string) error {
	devicePayload := make(payload.Payload).
		AddTrigger(languageGroupTitle, payload.TriggerDirectory, payload.TouchBar, payload.ActionTypeExecuteScript, false).
		AddIcon(payload.IconCharacter, payload.DefaultHeight, false)
	languageUUID, err := b.addTrigger(ctx, devicePayload, payload.SettingsOrderLanguage, parentUUID)
	if err != nil {
		return fmt.Errorf("cannot create language group: %w", err)
	}

	closePayload := make(payload.Payload).AddClose(settingsTitle)
	if _, err = b.addTrigger(ctx, closePayload, payload.LanguageOrderCloseGroup, languageUUID); err != nil {
		return fmt.Errorf("cannot create close trigger: %w", err)
	}

	rendered, err := b.renderer.Render("print_selected_language", map[string]any{"AppAddress": b.appAddress})
	if err != nil {
		return fmt.Errorf("cannot render print_selected_language script: %w", err)
	}
	selectedDevicePayload := make(payload.Payload).
		AddTrigger(selectedLanguageTitle, payload.TriggerShellScript, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
		AddShell(rendered, payload.IntervalDefault, payload.ShellTypeNone).
		AddMap(map[string]any{"BTTTriggerConfig": map[string]any{"BTTTouchBarFreeSpaceAfterButton": 25}})
	if _, err = b.addTrigger(ctx, selectedDevicePayload, payload.LanguageOrderSelectedLanguage, languageUUID); err != nil {
		return fmt.Errorf("cannot create selected language trigger: %w", err)
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
			AddShell(rendered, payload.IntervalNone, payload.ShellTypeAdditional)
		if _, err = b.addTrigger(ctx, languagePayload, payload.LanguageOrderSelectedLanguage+payload.Order(1+i), languageUUID); err != nil {
			return fmt.Errorf("cannot create select language trigger: %w", err)
		}
	}

	return nil
}

func (b *btt) addViewModeSection(ctx context.Context, parentUUID string) error {
	viewModePayload := make(payload.Payload).
		AddTrigger(viewModeTitle, payload.TriggerDirectory, payload.TouchBar, payload.ActionTypeExecuteScript, false).
		AddIcon(payload.IconMacwindow, payload.DefaultHeight, false)
	viewModeUUID, err := b.addTrigger(ctx, viewModePayload, payload.SettingsOrderViewMode, parentUUID)
	if err != nil {
		return fmt.Errorf("cannot create view mode group: %w", err)
	}

	closePayload := make(payload.Payload).AddClose(settingsTitle)
	if _, err = b.addTrigger(ctx, closePayload, payload.ViewModeOrderCloseGroup, viewModeUUID); err != nil {
		return fmt.Errorf("cannot create close trigger: %w", err)
	}

	var rendered string
	if rendered, err = b.renderer.Render("print_selected_view_mode", map[string]any{"AppAddress": b.appAddress}); err != nil {
		return fmt.Errorf("cannot render print_selected_view_mode script: %w", err)
	}
	selectedViewModePayload := make(payload.Payload).
		AddTrigger(selectedViewModeTitle, payload.TriggerShellScript, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
		AddShell(rendered, payload.IntervalDefault, payload.ShellTypeNone).
		AddMap(map[string]any{"BTTTriggerConfig": map[string]any{"BTTTouchBarFreeSpaceAfterButton": 25}})
	if _, err = b.addTrigger(ctx, selectedViewModePayload, payload.ViewModeOrderSelectedViewMode, viewModeUUID); err != nil {
		return fmt.Errorf("cannot create selected view mode trigger: %w", err)
	}

	viewModes := []string{ViewModeEmbed, ViewModeClean}
	for i, viewMode := range viewModes {
		rendered, err = b.renderer.Render(
			"select_view_mode",
			map[string]any{"AppAddress": b.appAddress, "ViewMode": viewMode},
		)
		if err != nil {
			return fmt.Errorf("cannot render select_view_mode script: %w", err)
		}

		viewModeOptionPayload := make(payload.Payload).
			AddTrigger(viewMode, payload.TriggerTouchBarButton, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
			AddShell(rendered, payload.IntervalNone, payload.ShellTypeAdditional)
		if _, err = b.addTrigger(ctx, viewModeOptionPayload, payload.ViewModeOrderSelectedViewMode+payload.Order(1+i), viewModeUUID); err != nil {
			return fmt.Errorf("cannot create view mode option trigger: %w", err)
		}
	}

	return nil
}

func (b *btt) addFloatingSection(ctx context.Context, parentUUID string) error {
	// Add floating menu
	floatingMenuPayload := make(
		payload.Payload,
	).AddFloatingMenu(floatingStateGroupTitle, payload.TriggerFloatingMenu, payload.FloatingMenu, true)
	floatingMenuUUID, err := b.addTrigger(ctx, floatingMenuPayload, 0, "")
	if err != nil {
		return fmt.Errorf("cannot create floating menu: %w", err)
	}
	rendered, err := b.renderer.Render("floating_page", map[string]any{"AppAddress": b.appAddress})
	if err != nil {
		return fmt.Errorf("cannot render floating_page html: %w", err)
	}

	webViewPayload := make(
		payload.Payload,
	).AddFloatingMenu(streamingTextTitle, payload.TriggerWebView, payload.FloatingMenu, false).
		AddMap(map[string]any{
			"BTTMenuConfig": map[string]any{
				"BTTMenuItemText": rendered,
			},
		})

	if _, err = b.addTrigger(ctx, webViewPayload, 0, floatingMenuUUID); err != nil {
		return fmt.Errorf("cannot create web view : %w", err)
	}

	// Add floating section in settings
	floatingPayload := make(payload.Payload).
		AddTrigger(floatingStateGroupTitle, payload.TriggerDirectory, payload.TouchBar, payload.ActionTypeExecuteScript, false).
		AddIcon(payload.IconMacwindow, payload.DefaultHeight, false)
	floatingUUID, err := b.addTrigger(ctx, floatingPayload, payload.SettingsOrderFloating, parentUUID)
	if err != nil {
		return fmt.Errorf("cannot create floating group: %w", err)
	}

	closePayload := make(payload.Payload).AddClose(settingsTitle)
	if _, err = b.addTrigger(ctx, closePayload, payload.FloatingOrderCloseGroup, floatingUUID); err != nil {
		return fmt.Errorf("cannot create close trigger: %w", err)
	}

	if rendered, err = b.renderer.Render("print_selected_floating_state", map[string]any{"AppAddress": b.appAddress}); err != nil {
		return fmt.Errorf("cannot render print_selected_floating_state script: %w", err)
	}
	selectedFloatingPayload := make(payload.Payload).
		AddTrigger(selectedFloatingStateTitle, payload.TriggerShellScript, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
		AddShell(rendered, payload.IntervalDefault, payload.ShellTypeNone).
		AddMap(map[string]any{"BTTTriggerConfig": map[string]any{"BTTTouchBarFreeSpaceAfterButton": 25}})
	if _, err = b.addTrigger(ctx, selectedFloatingPayload, payload.FloatingOrderSelectedState, floatingUUID); err != nil {
		return fmt.Errorf("cannot create selected floating trigger: %w", err)
	}

	floatingStates := []string{"Shown", "Hidden"}
	for i, floatingState := range floatingStates {
		rendered, err = b.renderer.Render(
			"select_floating_state",
			map[string]any{"AppAddress": b.appAddress, "FloatingState": floatingState},
		)
		if err != nil {
			return fmt.Errorf("cannot render select_floating_state script: %w", err)
		}

		statePayload := make(payload.Payload).
			AddTrigger(floatingState, payload.TriggerTouchBarButton, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
			AddShell(rendered, payload.IntervalNone, payload.ShellTypeAdditional)
		if _, err = b.addTrigger(ctx, statePayload, payload.FloatingOrderSelectedState+payload.Order(1+i), floatingUUID); err != nil {
			return fmt.Errorf("cannot create floating state trigger: %w", err)
		}
	}

	return nil
}

func (b *btt) addMetricsSection(ctx context.Context, parentUUID string) error {
	metricsPayload := make(payload.Payload).
		AddTrigger(metricsGroupTitle, payload.TriggerDirectory, payload.TouchBar, payload.ActionTypeExecuteScript, false).
		AddIcon(payload.IconChartBar, payload.DefaultHeight, false)
	metricsUUID, err := b.addTrigger(ctx, metricsPayload, payload.SettingsOrderMetrics, parentUUID)
	if err != nil {
		return fmt.Errorf("cannot create metrics group: %w", err)
	}

	closePayload := make(payload.Payload).AddClose(settingsTitle)
	if _, err = b.addTrigger(ctx, closePayload, payload.MetricsOrderCloseGroup, metricsUUID); err != nil {
		return fmt.Errorf("cannot create close trigger: %w", err)
	}

	for _, s := range []struct {
		template string
		metric   string
		title    string
		order    payload.Order
		icon     string
	}{
		{
			"print_size_metric",
			metrics.MetricBytesReadFromAudio,
			"Read",
			payload.MetricsOrderReadAudio,
			payload.IconWaveform,
		},
		{
			"print_size_metric",
			metrics.MetricBytesSentToGoogleSpeech,
			"Sent",
			payload.MetricsOrderSentAudio,
			payload.IconSquareAndArrowUp,
		},
		{
			"print_duration_metric",
			metrics.MetricMillisecondsSentToGoogleSpeech,
			"Sent",
			payload.MetricsOrderSendAudioMs,
			payload.IconSquareAndArrowUpBadgeClock,
		},
		{
			"print_size_metric",
			metrics.MetricBytesWrittenOnDisk,
			"Burnt",
			payload.MetricsOrderBurntAudio,
			payload.IconFlame,
		},
		{
			"print_raw",
			metrics.MetricConnectsToGoogleSpeech,
			"Connections",
			payload.MetricsOrderConnections,
			payload.IconAppConnectedToAppBelowFill,
		},
		{
			"print_raw",
			metrics.MetricTotalRunningTasks,
			"Tasks",
			payload.MetricsOrderTasks,
			payload.IconListBullet,
		},

		{
			"print_raw",
			metrics.MetricTotalOpenSockets,
			"Sockets",
			payload.MetricsOrderSockets,
			payload.IconRectangleStack,
		},
	} {
		var rendered string
		if rendered, err = b.renderer.Render(s.template, map[string]any{"AppAddress": b.appAddress, "Metric": s.metric, "Title": s.title}); err != nil {
			return fmt.Errorf("cannot render %s script: %w", s.template, err)
		}
		metricPayload := make(payload.Payload).
			AddTrigger(s.title, payload.TriggerShellScript, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
			AddShell(rendered, payload.IntervalMetrics, payload.ShellTypeNone).
			AddIcon(s.icon, payload.DefaultHeight, false)
		if _, err = b.addTrigger(ctx, metricPayload, s.order, metricsUUID); err != nil {
			return fmt.Errorf("cannot create metric trigger: %w", err)
		}
	}

	return nil
}

func (b *btt) addCleanModeViewSection(ctx context.Context) error {
	cleanModeViewPayload := make(
		payload.Payload,
	).AddTrigger(cleanViewTitle, payload.TriggerDirectory, payload.TouchBar, payload.ActionTypeNone, true)

	cleanModeViewUUID, err := b.addTrigger(ctx, cleanModeViewPayload, payload.MainOrderViewMode, "")
	if err != nil {
		return fmt.Errorf("cannot create clean mode view group: %w", err)
	}

	return b.addMainAppTrigger(ctx, cleanViewAppTitle, cleanModeViewUUID)
}

func (b *btt) addMainSection(ctx context.Context) error {
	jsonPayload := map[string]any{
		"BTTPredefinedActionType": payload.ActionTypeOpenGroup,
		"BTTOpenGroupWithName":    settingsTitle,
	}
	encoded, err := encodeForScript(jsonPayload)
	if err != nil {
		return err
	}

	rendered, err := b.renderer.Render(
		"open_settings",
		map[string]any{"AppAddress": b.appAddress, "BttAddress": b.bttAddress, "Query": encoded},
	)
	if err != nil {
		return fmt.Errorf("cannot render print_selected_language script: %w", err)
	}

	namedPayload := make(payload.Payload).
		AddTrigger(settingsTitle, payload.TriggerNamed, payload.OtherTriggers, payload.ActionTypeEmptyPlaceholder, false).
		AddShell(rendered, payload.IntervalNone, payload.ShellTypeAdditional)
	if _, err = b.addTrigger(ctx, namedPayload, payload.NamedOrderSettings, ""); err != nil {
		return fmt.Errorf("cannot create named trigger: %w", err)
	}

	return b.addMainAppTrigger(ctx, appTitle, "")
}

func (b *btt) addMainAppTrigger(ctx context.Context, name, parentUUID string) error {
	toggleRendered, err := b.renderer.Render("toggle", map[string]any{"AppAddress": b.appAddress})
	if err != nil {
		return fmt.Errorf("cannot render toggle script: %w", err)
	}
	listenSocketRendered, err := b.renderer.Render("listen_socket", map[string]any{})
	if err != nil {
		return fmt.Errorf("cannot render listen_socket script: %w", err)
	}

	mainPayload := make(payload.Payload).
		AddTrigger(name, payload.TriggerShellScript, payload.TouchBar, payload.ActionTypeEmptyPlaceholder, false).
		AddShell(toggleRendered, payload.IntervalNone, payload.ShellTypeEmbed).
		AddShell(listenSocketRendered, payload.IntervalNone, payload.ShellTypeNone).
		AddMap(map[string]any{
			"BTTTriggerConfig": map[string]any{
				"BTTTouchBarLongPressActionName": settingsTitle,
				"BTTTouchBarButtonFontSize":      12,
				"BTTTouchBarButtonColor":         "0.0, 0.0, 0.0, 255.0",
				"BTTTouchBarButtonTextAlignment": 0,
			},
		})
	if _, err = b.addTrigger(ctx, mainPayload, payload.MainOrderSubs, parentUUID); err != nil {
		return fmt.Errorf("cannot create main trigger: %w", err)
	}

	return nil
}

func encodeForScript(payload map[string]any) (string, error) {
	// TODO: code is repeated
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return "", fmt.Errorf("cannot encode json payload: %w", err)
	}

	query := url.Values{}
	query.Set("json", buf.String())
	encoded := strings.ReplaceAll(query.Encode(), "+", "%20")

	return encoded, nil
}
