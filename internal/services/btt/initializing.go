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

func (i *initializingComponent) Initialize(ctx context.Context) error {
	after := trigger.NewTrigger().AddOrder(trigger.OrderSettings)
	settingsDir, err := i.addSettingsSection(ctx, "", after)
	if err != nil {
		return fmt.Errorf("cannot add settings section: %w", err)
	}
	cleanViewDir, err := i.addCleanViewSection(ctx, "", settingsDir)
	if err != nil {
		return fmt.Errorf("cannot add clean view section: %w", err)
	}
	if err = i.addAppSection(ctx, "", cleanViewDir); err != nil {
		return fmt.Errorf("cannot add app section: %w", err)
	}

	return nil
}

func (i *initializingComponent) addSettingsSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	settingsDir := trigger.NewHiddenDir(trigger.TitleSettingsDir).AddOrderAfter(after)
	settingsUUID, err := i.client.AddTrigger(ctx, settingsDir, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot add settings directory: %w", err)
	}

	return settingsDir, i.addTriggers(ctx, settingsUUID, []addingTriggerFunc{
		i.addCloseSettingsSection,
		i.addStatusSection,
		i.addDeviceSection,
		i.addLanguageSection,
		i.addViewModeSection,
		i.addFloatingStateSection,
		i.addClipboardSection,
		i.addMetricsSection,
	}, nil)
}

func (i *initializingComponent) addCloseSettingsSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	rendered := i.renderer.CloseSettings(
		string(ViewModeClean),
		trigger.NewCloseDirAction(),
		trigger.NewOpenDirAction(trigger.TitleCleanViewDir),
	)
	closeDir := trigger.NewTapButton(trigger.TitleCloseDir, rendered).
		AddCloseIcon().
		AddOrderAfter(after)
	if _, err := i.client.AddTrigger(ctx, closeDir, parentUUID); err != nil {
		return nil, fmt.Errorf("cannot add close settings trigger: %w", err)
	}

	return closeDir, nil
}

func (i *initializingComponent) addStatusSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	printStatus := trigger.NewStatusInfoButton("⏳", i.renderer.PrintStatus()).AddOrderAfter(after)
	if _, err := i.client.AddTrigger(ctx, printStatus, parentUUID); err != nil {
		return nil, fmt.Errorf("cannot add print status trigger: %w", err)
	}

	return printStatus, nil
}

func (i *initializingComponent) addDeviceSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	deviceDir := trigger.NewDirButton(trigger.TitleDeviceDir, trigger.IconMicrophone).AddOrderAfter(after)
	deviceUUID, err := i.client.AddTrigger(ctx, deviceDir, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot add device directory: %w", err)
	}

	closeDir := trigger.NewOpenDirButton(trigger.TitleSettingsDir)
	selectedDevice := trigger.NewSettingsInfoButton(trigger.TitleSelectedDevice, i.renderer.PrintSelectedDevice())

	return deviceDir, i.addTriggers(ctx, deviceUUID, []addingTriggerFunc{
		i.wrapAddingTrigger(closeDir),
		i.wrapAddingTrigger(selectedDevice),
	}, nil)
}

func (i *initializingComponent) addLanguageSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	languageDir := trigger.NewDirButton(trigger.TitleLanguageDir, trigger.IconCharacter).AddOrderAfter(after)
	languageUUID, err := i.client.AddTrigger(ctx, languageDir, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot add language directory: %w", err)
	}

	closeDir := trigger.NewOpenDirButton(trigger.TitleSettingsDir)
	selectedLanguage := trigger.NewSettingsInfoButton(trigger.TitleSelectedLanguage, i.renderer.PrintSelectedLanguage())

	addingTriggers := []addingTriggerFunc{
		i.wrapAddingTrigger(closeDir),
		i.wrapAddingTrigger(selectedLanguage),
	}
	for _, l := range i.languages {
		languageButton := trigger.NewTapButton(trigger.Title(l), i.renderer.SelectLanguage(l))
		addingTriggers = append(addingTriggers, i.wrapAddingTrigger(languageButton))
	}

	return languageDir, i.addTriggers(ctx, languageUUID, addingTriggers, nil)
}

//nolint:dupl
func (i *initializingComponent) addViewModeSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	viewModeDir := trigger.NewDirButton(trigger.TitleViewModeDir, trigger.IconMacwindow).AddOrderAfter(after)
	viewModeUUID, err := i.client.AddTrigger(ctx, viewModeDir, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot add view mode directory: %w", err)
	}

	closeDir := trigger.NewOpenDirButton(trigger.TitleSettingsDir)
	selectedViewMode := trigger.NewSettingsInfoButton(trigger.TitleSelectedViewMode, i.renderer.PrintSelectedViewMode())
	embedButton := trigger.NewTapButton(trigger.Title(ViewModeEmbed), i.renderer.SelectViewMode(string(ViewModeEmbed)))
	cleanButton := trigger.NewTapButton(trigger.Title(ViewModeClean), i.renderer.SelectViewMode(string(ViewModeClean)))

	return viewModeDir, i.addTriggers(ctx, viewModeUUID, []addingTriggerFunc{
		i.wrapAddingTrigger(closeDir),
		i.wrapAddingTrigger(selectedViewMode),
		i.wrapAddingTrigger(embedButton),
		i.wrapAddingTrigger(cleanButton),
	}, nil)
}

func (i *initializingComponent) addFloatingStateSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	// Adding floating menu
	floatingMenu := trigger.NewFloatingMenu(trigger.TitleStreamingText)
	floatingMenuUUID, err := i.client.AddTrigger(ctx, floatingMenu, "")
	if err != nil {
		return nil, fmt.Errorf("cannot add floating menu: %w", err)
	}
	webView := trigger.NewWebView(trigger.TitleStreamingText, i.renderer.FloatingPage())
	if _, err = i.client.AddTrigger(ctx, webView, floatingMenuUUID); err != nil {
		return nil, fmt.Errorf("cannot add floating view: %w", err)
	}

	// Adding buttons in settings section
	floatingStateDir := trigger.NewDirButton(trigger.TitleFloatingDir, trigger.IconMacwindow).AddOrderAfter(after)
	floatingStateUUID, err := i.client.AddTrigger(ctx, floatingStateDir, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot add floating state directory: %w", err)
	}

	closeDir := trigger.NewOpenDirButton(trigger.TitleSettingsDir)
	selectedFloatingState := trigger.NewSettingsInfoButton(
		trigger.TitleSelectedFloating,
		i.renderer.PrintSelectedFloating(),
	)
	shownButton := trigger.NewTapButton(trigger.Title(FloatingShown), i.renderer.SelectFloating(string(FloatingShown)))
	hiddenButton := trigger.NewTapButton(
		trigger.Title(FloatingHidden),
		i.renderer.SelectFloating(string(FloatingHidden)),
	)

	return floatingStateDir, i.addTriggers(ctx, floatingStateUUID, []addingTriggerFunc{
		i.wrapAddingTrigger(closeDir),
		i.wrapAddingTrigger(selectedFloatingState),
		i.wrapAddingTrigger(shownButton),
		i.wrapAddingTrigger(hiddenButton),
	}, nil)
}

//nolint:dupl
func (i *initializingComponent) addClipboardSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	clipboardDir := trigger.NewDirButton(trigger.TitleClipboardDir, trigger.IconClipboard).AddOrderAfter(after)
	clipboardUUID, err := i.client.AddTrigger(ctx, clipboardDir, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot add clipboard directory: %w", err)
	}

	closeDir := trigger.NewOpenDirButton(trigger.TitleSettingsDir)
	selectedClipboard := trigger.NewSettingsInfoButton(
		trigger.TitleSelectedClipboard,
		i.renderer.PrintSelectedClipboard(),
	)
	shownButton := trigger.NewTapButton(
		trigger.Title(ClipboardShown),
		i.renderer.SelectClipboard(string(ClipboardShown)),
	)
	hiddenButton := trigger.NewTapButton(
		trigger.Title(ClipboardHidden),
		i.renderer.SelectClipboard(string(ClipboardHidden)),
	)

	return clipboardDir, i.addTriggers(ctx, clipboardUUID, []addingTriggerFunc{
		i.wrapAddingTrigger(closeDir),
		i.wrapAddingTrigger(selectedClipboard),
		i.wrapAddingTrigger(shownButton),
		i.wrapAddingTrigger(hiddenButton),
	}, nil)
}

func (i *initializingComponent) addMetricsSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	metricsDir := trigger.NewDirButton(trigger.TitleMetricsDir, trigger.IconChartBar).AddOrderAfter(after)
	metricsUUID, err := i.client.AddTrigger(ctx, metricsDir, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot add metrics directory: %w", err)
	}

	closeDir := trigger.NewOpenDirButton(trigger.TitleSettingsDir)

	readScript := i.renderer.PrintMetric(tmpl.MetricTemplateSize, metrics.MetricBytesReadFromAudio, "Read")
	readButton := trigger.NewMetricsInfoButton("Read", readScript, trigger.IconWaveform)

	sentBytesScript := i.renderer.PrintMetric(tmpl.MetricTemplateSize, metrics.MetricBytesSentToGoogleSpeech, "Sent")
	sentBytesButton := trigger.NewMetricsInfoButton("Sent", sentBytesScript, trigger.IconSquareAndArrowUp)

	sentDurationScript := i.renderer.PrintMetric(
		tmpl.MetricTemplateDuration,
		metrics.MetricMillisecondsSentToGoogleSpeech,
		"Sent",
	)
	sentDurationButton := trigger.NewMetricsInfoButton(
		"Sent",
		sentDurationScript,
		trigger.IconSquareAndArrowUpBadgeClock,
	)

	burntScript := i.renderer.PrintMetric(tmpl.MetricTemplateSize, metrics.MetricBytesWrittenOnDisk, "Burnt")
	burntButton := trigger.NewMetricsInfoButton("Burnt", burntScript, trigger.IconFlame)

	connectionsScript := i.renderer.PrintMetric(
		tmpl.MetricTemplateRaw,
		metrics.MetricConnectsToGoogleSpeech,
		"Connections",
	)
	connectionsButton := trigger.NewMetricsInfoButton(
		"Connections",
		connectionsScript,
		trigger.IconAppConnectedToAppBelowFill,
	)

	tasksScript := i.renderer.PrintMetric(tmpl.MetricTemplateRaw, metrics.MetricTotalRunningTasks, "Tasks")
	tasksButton := trigger.NewMetricsInfoButton("Tasks", tasksScript, trigger.IconListBullet)

	socketsScript := i.renderer.PrintMetric(tmpl.MetricTemplateRaw, metrics.MetricTotalOpenSockets, "Sockets")
	socketsButton := trigger.NewMetricsInfoButton("Sockets", socketsScript, trigger.IconRectangleStack)

	return metricsDir, i.addTriggers(ctx, metricsUUID, []addingTriggerFunc{
		i.wrapAddingTrigger(closeDir),
		i.wrapAddingTrigger(readButton),
		i.wrapAddingTrigger(sentBytesButton),
		i.wrapAddingTrigger(sentDurationButton),
		i.wrapAddingTrigger(burntButton),
		i.wrapAddingTrigger(connectionsButton),
		i.wrapAddingTrigger(tasksButton),
		i.wrapAddingTrigger(socketsButton),
	}, nil)
}

func (i *initializingComponent) addCleanViewSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) (trigger.Trigger, error) {
	cleanViewDir := trigger.NewHiddenDir(trigger.TitleCleanViewDir).AddOrderAfter(after)
	cleanViewDirUUID, err := i.client.AddTrigger(ctx, cleanViewDir, parentUUID)
	if err != nil {
		return nil, fmt.Errorf("cannot add clean mode directory: %w", err)
	}

	return cleanViewDir, i.addAppTriggers(
		ctx,
		trigger.TitleCleanViewApp,
		trigger.TitleCleanViewClipboard,
		cleanViewDirUUID,
		nil,
	)
}

func (i *initializingComponent) addAppSection(
	ctx context.Context,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) error {
	// Create a named trigger to manage "open settings"
	openDirAction := trigger.NewOpenDirAction(trigger.TitleSettingsDir)
	openSettings := trigger.NewNamedTrigger(trigger.TitleOpenSettings, i.renderer.OpenSettings(openDirAction))
	if _, err := i.client.AddTrigger(ctx, openSettings, ""); err != nil {
		return fmt.Errorf("cannot add open settings named trigger: %w", err)
	}

	return i.addAppTriggers(ctx, trigger.TitleApp, trigger.TitleClipboard, parentUUID, after)
}

func (i *initializingComponent) addAppTriggers(
	ctx context.Context,
	appTitle trigger.Title,
	clipboardTitle trigger.Title,
	parentUUID trigger.UUID,
	after trigger.Trigger,
) error {
	clipboard := trigger.NewTapIconButton(clipboardTitle, i.renderer.CopyText(), trigger.IconClipboard).AddDisabled()
	app := trigger.NewInfoButton(appTitle, i.renderer.AppPlaceholder(), 0).
		AddTapScript(i.renderer.Toggle()).
		AddLongTapTrigger(trigger.TitleOpenSettings).
		AddReadableFormat().
		AddOrderAfter(after)

	return i.addTriggers(ctx, parentUUID, []addingTriggerFunc{
		i.wrapAddingTrigger(clipboard),
		i.wrapAddingTrigger(app),
	}, after)
}

func (i *initializingComponent) Clear(ctx context.Context) error {
	action := trigger.NewCloseDirAction()
	if err := i.client.TriggerAction(ctx, action); err != nil {
		return fmt.Errorf("cannot close directory: %w", err)
	}
	triggers, err := i.client.GetTriggers(ctx, "")
	if err != nil {
		return fmt.Errorf("cannot get triggers: %w", err)
	}

	if err = i.client.DeleteTriggers(ctx, triggers); err != nil {
		return fmt.Errorf("cannot delete triggers: %w", err)
	}

	return nil
}

func (i *initializingComponent) wrapAddingTrigger(t trigger.Trigger) addingTriggerFunc {
	return func(ctx context.Context, parentUUID trigger.UUID, after trigger.Trigger) (trigger.Trigger, error) {
		t.AddOrderAfter(after)
		if _, err := i.client.AddTrigger(ctx, t, parentUUID); err != nil {
			return nil, fmt.Errorf("cannot add %s: %w", t.ErrorContext(), err)
		}
		return t, nil
	}
}

func (i *initializingComponent) addTriggers(
	ctx context.Context,
	parentUUID trigger.UUID,
	fns []addingTriggerFunc,
	after trigger.Trigger,
) error {
	var err error
	for _, fn := range fns {
		after, err = fn(ctx, parentUUID, after)
		if err != nil {
			return err
		}
	}
	return nil
}
