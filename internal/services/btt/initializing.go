package btt

import (
	"context"
	"fmt"

	"live2text/internal/services/btt/client"
	"live2text/internal/services/btt/client/trigger"
	"live2text/internal/services/btt/tmpl"
	"live2text/internal/services/metrics"
)

type initializingComponent struct {
	client    client.Client
	renderer  tmpl.Renderer
	languages []string
}

func NewInitializingComponent(client client.Client, renderer tmpl.Renderer, languages []string) InitializingComponent {
	return &initializingComponent{client: client, renderer: renderer, languages: languages}
}

type addingTriggerFunc func(ctx context.Context, parentUUID trigger.UUID, after trigger.Trigger) (trigger.Trigger, error)

func (b *initializingComponent) Initialize(ctx context.Context) error {
	after := trigger.NewTrigger().AddOrder(trigger.OrderSettings)
	settingsDir, err := b.addSettingsSection(ctx, "", after)
	if err != nil {
		return fmt.Errorf("cannot add settings section: %w", err)
	}
	cleanViewDir, err := b.addCleanViewSection(ctx, "", settingsDir)
	if err != nil {
		return fmt.Errorf("cannot add clean view section: %w", err)
	}
	_, err = b.addAppSection(ctx, "", cleanViewDir)
	if err != nil {
		return fmt.Errorf("cannot add app section: %w", err)
	}

	return nil
}

func (b *initializingComponent) addSettingsSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	settingsDir := trigger.NewHiddenDir(trigger.TitleSettingsDir).AddOrderAfter(after)
	settingsUUID, err := b.client.AddTrigger(ctx, settingsDir, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot add settings directory: %w", err)
	}

	return settingsDir, b.addTriggers(ctx, settingsUUID, []addingTriggerFunc{
		b.addCloseSettingsSection,
		b.addStatusSection,
		b.addDeviceSection,
		b.addLanguageSection,
		b.addViewModeSection,
		b.addFloatingStateSection,
		b.addMetricsSection,
	})
}

func (b *initializingComponent) addCloseSettingsSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	rendered := b.renderer.CloseSettings(
		string(ViewModeClean),
		trigger.NewCloseDirAction(),
		trigger.NewOpenDirAction(trigger.TitleCleanViewDir),
	)
	closeDir := trigger.NewTapButton(trigger.TitleCloseDir, rendered).
		AddCloseIcon().
		AddOrderAfter(after)
	if _, err := b.client.AddTrigger(ctx, closeDir, parentUUID); err != nil {
		return nil, fmt.Errorf("cannot add close settings trigger: %w", err)
	}

	return closeDir, nil
}

func (b *initializingComponent) addStatusSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	printStatus := trigger.NewStatusInfoButton("⏳", b.renderer.PrintStatus()).AddOrderAfter(after)
	if _, err := b.client.AddTrigger(ctx, printStatus, parentUUID); err != nil {
		return nil, fmt.Errorf("cannot add print status trigger: %w", err)
	}

	return printStatus, nil
}

func (b *initializingComponent) addDeviceSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	deviceDir := trigger.NewDirButton(trigger.TitleDeviceDir, trigger.IconMicrophone).AddOrderAfter(after)
	deviceUUID, err := b.client.AddTrigger(ctx, deviceDir, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot add device directory: %w", err)
	}

	closeDir := trigger.NewOpenDirButton(trigger.TitleSettingsDir)
	selectedDevice := trigger.NewSettingsInfoButton(trigger.TitleSelectedDevice, b.renderer.PrintSelectedDevice())

	return deviceDir, b.addTriggers(ctx, deviceUUID, []addingTriggerFunc{
		b.wrapAddingTrigger(closeDir),
		b.wrapAddingTrigger(selectedDevice),
	})
}

func (b *initializingComponent) addLanguageSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	languageDir := trigger.NewDirButton(trigger.TitleLanguageDir, trigger.IconCharacter).AddOrderAfter(after)
	languageUUID, err := b.client.AddTrigger(ctx, languageDir, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot add language directory: %w", err)
	}

	closeDir := trigger.NewOpenDirButton(trigger.TitleSettingsDir)
	selectedLanguage := trigger.NewSettingsInfoButton(trigger.TitleSelectedLanguage, b.renderer.PrintSelectedLanguage())

	addingTriggers := []addingTriggerFunc{
		b.wrapAddingTrigger(closeDir),
		b.wrapAddingTrigger(selectedLanguage),
	}
	for _, l := range b.languages {
		languageButton := trigger.NewTapButton(trigger.Title(l), b.renderer.SelectLanguage(l))
		addingTriggers = append(addingTriggers, b.wrapAddingTrigger(languageButton))
	}

	return languageDir, b.addTriggers(ctx, languageUUID, addingTriggers)
}

func (b *initializingComponent) addViewModeSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	viewModeDir := trigger.NewDirButton(trigger.TitleViewModeDir, trigger.IconMacwindow).AddOrderAfter(after)
	viewModeUUID, err := b.client.AddTrigger(ctx, viewModeDir, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot add view mode directory: %w", err)
	}

	closeDir := trigger.NewOpenDirButton(trigger.TitleSettingsDir)
	selectedViewMode := trigger.NewSettingsInfoButton(trigger.TitleSelectedViewMode, b.renderer.PrintSelectedViewMode())
	embedButton := trigger.NewTapButton(trigger.Title(ViewModeEmbed), b.renderer.SelectViewMode(string(ViewModeEmbed)))
	cleanButton := trigger.NewTapButton(trigger.Title(ViewModeClean), b.renderer.SelectViewMode(string(ViewModeClean)))

	return viewModeDir, b.addTriggers(ctx, viewModeUUID, []addingTriggerFunc{
		b.wrapAddingTrigger(closeDir),
		b.wrapAddingTrigger(selectedViewMode),
		b.wrapAddingTrigger(embedButton),
		b.wrapAddingTrigger(cleanButton),
	})
}

func (b *initializingComponent) addFloatingStateSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	// Adding floating menu
	floatingMenu := trigger.NewFloatingMenu(trigger.TitleStreamingText)
	floatingMenuUUID, err := b.client.AddTrigger(ctx, floatingMenu, "")
	if err != nil {
		return nil, fmt.Errorf("cannot add floating menu: %w", err)
	}
	webView := trigger.NewWebView(trigger.TitleStreamingText, b.renderer.FloatingPage())
	if _, err = b.client.AddTrigger(ctx, webView, floatingMenuUUID); err != nil {
		return nil, fmt.Errorf("cannot add floating view: %w", err)
	}

	// Adding buttons in settings section
	floatingStateDir := trigger.NewDirButton(trigger.TitleFloatingDir, trigger.IconMacwindow).AddOrderAfter(after)
	floatingStateUUID, err := b.client.AddTrigger(ctx, floatingStateDir, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot add floating state directory: %w", err)
	}

	closeDir := trigger.NewOpenDirButton(trigger.TitleSettingsDir)
	selectedFloatingState := trigger.NewSettingsInfoButton(
		trigger.TitleSelectedFloating,
		b.renderer.PrintSelectedFloatingState(),
	)
	shownButton := trigger.NewTapButton(FloatingShown, b.renderer.SelectFloatingState(FloatingShown))
	hiddenButton := trigger.NewTapButton(FloatingHidden, b.renderer.SelectFloatingState(FloatingHidden))

	return floatingStateDir, b.addTriggers(ctx, floatingStateUUID, []addingTriggerFunc{
		b.wrapAddingTrigger(closeDir),
		b.wrapAddingTrigger(selectedFloatingState),
		b.wrapAddingTrigger(shownButton),
		b.wrapAddingTrigger(hiddenButton),
	})
}

func (b *initializingComponent) addMetricsSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	metricsDir := trigger.NewDirButton(trigger.TitleMetricsDir, trigger.IconChartBar).AddOrderAfter(after)
	metricsUUID, err := b.client.AddTrigger(ctx, metricsDir, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot add metrics directory: %w", err)
	}

	closeDir := trigger.NewOpenDirButton(trigger.TitleSettingsDir)

	readScript := b.renderer.PrintMetric(tmpl.MetricTemplateSize, metrics.MetricBytesReadFromAudio, "Read")
	readButton := trigger.NewMetricsInfoButton("Read", readScript, trigger.IconWaveform)

	sentBytesScript := b.renderer.PrintMetric(tmpl.MetricTemplateSize, metrics.MetricBytesSentToGoogleSpeech, "Sent")
	sentBytesButton := trigger.NewMetricsInfoButton("Sent", sentBytesScript, trigger.IconSquareAndArrowUp)

	sentDurationScript := b.renderer.PrintMetric(
		tmpl.MetricTemplateDuration,
		metrics.MetricMillisecondsSentToGoogleSpeech,
		"Sent",
	)
	sentDurationButton := trigger.NewMetricsInfoButton(
		"Sent",
		sentDurationScript,
		trigger.IconSquareAndArrowUpBadgeClock,
	)

	burntScript := b.renderer.PrintMetric(tmpl.MetricTemplateSize, metrics.MetricBytesWrittenOnDisk, "Burnt")
	burntButton := trigger.NewMetricsInfoButton("Burnt", burntScript, trigger.IconFlame)

	connectionsScript := b.renderer.PrintMetric(
		tmpl.MetricTemplateRaw,
		metrics.MetricConnectsToGoogleSpeech,
		"Connections",
	)
	connectionsButton := trigger.NewMetricsInfoButton(
		"Connections",
		connectionsScript,
		trigger.IconAppConnectedToAppBelowFill,
	)

	tasksScript := b.renderer.PrintMetric(tmpl.MetricTemplateRaw, metrics.MetricTotalRunningTasks, "Tasks")
	tasksButton := trigger.NewMetricsInfoButton("Tasks", tasksScript, trigger.IconListBullet)

	socketsScript := b.renderer.PrintMetric(tmpl.MetricTemplateRaw, metrics.MetricTotalOpenSockets, "Sockets")
	socketsButton := trigger.NewMetricsInfoButton("Sockets", socketsScript, trigger.IconRectangleStack)

	return metricsDir, b.addTriggers(ctx, metricsUUID, []addingTriggerFunc{
		b.wrapAddingTrigger(closeDir),
		b.wrapAddingTrigger(readButton),
		b.wrapAddingTrigger(sentBytesButton),
		b.wrapAddingTrigger(sentDurationButton),
		b.wrapAddingTrigger(burntButton),
		b.wrapAddingTrigger(connectionsButton),
		b.wrapAddingTrigger(tasksButton),
		b.wrapAddingTrigger(socketsButton),
	})
}

func (b *initializingComponent) addTriggers(
	ctx context.Context,
	parentUUID trigger.UUID,
	fns []addingTriggerFunc,
) error {
	var after trigger.Trigger
	var err error
	for _, fn := range fns {
		after, err = fn(ctx, parentUUID, after)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *initializingComponent) wrapAddingTrigger(t trigger.Trigger) addingTriggerFunc {
	return func(ctx context.Context, parentUUID trigger.UUID, after trigger.Trigger) (trigger.Trigger, error) {
		t.AddOrderAfter(after)
		if _, err := b.client.AddTrigger(ctx, t, parentUUID); err != nil {
			return nil, fmt.Errorf("cannot add %s: %w", t.ErrorContext(), err)
		}
		return t, nil
	}
}

func (b *initializingComponent) addCleanViewSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	cleanViewDir := trigger.NewHiddenDir(trigger.TitleCleanViewDir).AddOrderAfter(after)
	cleanViewDirUUID, err := b.client.AddTrigger(ctx, cleanViewDir, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot add clean mode directory: %w", err)
	}

	_, err = b.addAppTrigger(ctx, trigger.TitleCleanViewApp, cleanViewDirUUID, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot add app trigger: %w", err)
	}

	return cleanViewDir, nil
}

func (b *initializingComponent) addAppSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	// Create a named trigger to manage "open settings"
	openDirAction := trigger.NewOpenDirAction(trigger.TitleSettingsDir)
	openSettings := trigger.NewNamedTrigger(trigger.TitleOpenSettings, b.renderer.OpenSettings(openDirAction))
	if _, err := b.client.AddTrigger(ctx, openSettings, parentUUID); err != nil {
		return nil, fmt.Errorf("cannot add open settings trigger: %w", err)
	}

	return b.addAppTrigger(ctx, trigger.TitleApp, parentUUID, after)
}

func (b *initializingComponent) addAppTrigger(
	ctx context.Context,
	title trigger.Title,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	app := trigger.NewInfoButton(title, b.renderer.AppPlaceholder(), 0).
		AddTapScript(b.renderer.Toggle()).
		AddLongTapTrigger(trigger.TitleOpenSettings).
		AddReadableFormat().
		AddOrderAfter(after)

	if _, err := b.client.AddTrigger(ctx, app, parentUUID); err != nil {
		return nil, fmt.Errorf("cannot add app trigger: %w", err)
	}

	return app, nil
}

func (b *initializingComponent) Clear(ctx context.Context) error {
	action := trigger.NewCloseDirAction()
	if err := b.client.TriggerAction(ctx, action); err != nil {
		return fmt.Errorf("cannot close directory: %w", err)
	}
	triggers, err := b.client.GetTriggers(ctx, "")
	if err != nil {
		return fmt.Errorf("cannot get triggers: %w", err)
	}

	if err = b.client.DeleteTriggers(ctx, triggers); err != nil {
		return fmt.Errorf("cannot delete triggers: %w", err)
	}

	return nil
}
