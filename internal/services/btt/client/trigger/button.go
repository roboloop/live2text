package trigger

// NewTapButton creates a new button with a shell script on a tap.
func NewTapButton(title Title, script string) Trigger {
	return newTrigger().
		init(title, triggerTouchBarButton, typeTouchBar).
		addAction(actionTypeEmptyPlaceholder).
		addTapScript(script)
}

// NewInfoButton creates a new button with cycled execution of a shell script on view.
func NewInfoButton(title Title, script string, interval float64) Trigger {
	return newTrigger().
		init(title, triggerShellScript, typeTouchBar).
		addAction(actionTypeEmptyPlaceholder).
		addCycledScript(script, interval)
}

// NewDirButton create a new button that opens a directory.
func NewDirButton(title Title, icon Icon) Trigger {
	return newTrigger().
		init(title, triggerDirectory, typeTouchBar).
		addIcon(icon, defaultHeight, false)
}

// NewStatusInfoButton creates a button with cycled execution of a shell script on view
// It's used for showing the state of an app.
func NewStatusInfoButton(title Title, script string) Trigger {
	return NewInfoButton(title, script, settingsInterval)
}

// NewSettingsInfoButton creates a new button with cycled execution of a shell script on view
// It's used for showing piece of settings.
func NewSettingsInfoButton(title Title, script string) Trigger {
	return NewInfoButton(title, script, settingsInterval).addFreeSpaceAfter()
}

// NewMetricsInfoButton creates a new button with cycled execution of grabbing metrics.
func NewMetricsInfoButton(title Title, script string, icon Icon) Trigger {
	return NewInfoButton(title, script, metricsInterval).addIcon(icon, defaultHeight, false)
}

// NewHiddenDir creates a new hidden directory.
func NewHiddenDir(title Title) Trigger {
	return newTrigger().init(title, triggerDirectory, typeTouchBar).addHidden()
}

// NewCloseDirButton creates a button that closes any open directories.
func NewCloseDirButton() Trigger {
	return newTrigger().
		init(TitleCloseDir, triggerTouchBarButton, typeTouchBar).
		addCloseDirAction().
		addCloseIcon()
}

// NewOpenDirButton creates a new button that opens another directory.
func NewOpenDirButton(dirName Title) Trigger {
	return newTrigger().
		init(TitleCloseDir, triggerTouchBarButton, typeTouchBar).
		addOpenDirAction(dirName).
		addCloseIcon()
}

// NewNamedTrigger creates a named trigger with a shell script on a long tap.
func NewNamedTrigger(title Title, script string) Trigger {
	return newTrigger().
		init(title, triggerNamed, typeOtherTriggers).
		addAction(actionTypeEmptyPlaceholder).
		addTapScript(script)
}

func (t Trigger) AddCloseIcon() Trigger {
	return t.addCloseIcon()
}

func (t Trigger) addCloseIcon() Trigger {
	return t.addIcon(IconXmarkCircleFill, 25, true)
}

func (t Trigger) addHidden() Trigger {
	t.addExtra(map[string]any{
		"BTTTriggerConfig": map[string]any{
			"BTTTouchBarButtonWidth":         0,
			"BTTTouchBarButtonUseFixedWidth": 1,
		},
	})

	return t
}

const defaultFreeSpaceAfter float64 = 25

func (t Trigger) addFreeSpaceAfter() Trigger {
	return t.addExtra(map[string]any{
		"BTTTriggerConfig": map[string]any{
			"BTTTouchBarFreeSpaceAfterButton": defaultFreeSpaceAfter,
		},
	})
}
