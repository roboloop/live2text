package trigger

type id float64

func (i id) float() float64 {
	return float64(i)
}

const (
	triggerDirectory      id = 630
	triggerTouchBarButton id = 629
	triggerNamed          id = 643
	triggerShellScript    id = 642
	triggerFloatingMenu   id = 767
	triggerWebView        id = 778
)

type bttType string

func (bt bttType) String() string {
	return string(bt)
}

const (
	typeTouchBar      bttType = "BTTTriggerTypeTouchBar"
	typeOtherTriggers bttType = "BTTTriggerTypeOtherTriggers"
	typeFloatingMenu  bttType = "BTTTriggerTypeFloatingMenu"
)

type actionType float64

func (at actionType) Float() float64 {
	return float64(at)
}

func (at actionType) String() string {
	switch at {
	case actionTypeEmptyPlaceholder:
		return "Empty Placeholder"
	case actionTypeExecuteScript:
		return "Execute Script"
	case actionTypeCloseGroup:
		return "Close Group"
	case actionTypeOpenGroup:
		return "Open Group"
	case actionTypeOpenFloatingMenu:
		return "Open Floating Menu"
	case actionTypeCloseFloatingMenu:
		return "Close Floating Menu"
	case actionTypeToggleFloatingMenu:
		return "Toggle Floating Menu"
	}
	return "Unknown Action Type"
}

const (
	actionTypeEmptyPlaceholder   actionType = 366
	actionTypeExecuteScript      actionType = 206
	actionTypeCloseGroup         actionType = 191
	actionTypeOpenGroup          actionType = 205
	actionTypeOpenFloatingMenu   actionType = 386
	actionTypeCloseFloatingMenu  actionType = 387
	actionTypeToggleFloatingMenu actionType = 388
)

type Title string

func (t Title) String() string {
	return string(t)
}

const (
	TitleApp                Title = "Live2Text App"
	TitleClipboard          Title = "Copy Text"
	TitleCleanViewApp       Title = "Clean View App"
	TitleCleanViewClipboard Title = "Clean View Copy Text"

	TitleSettingsDir      Title = "Settings"
	TitleCleanViewDir     Title = "Clean View"
	TitleCloseSettingsDir Title = "Close Settings"
	TitleCloseDir         Title = "Close Directory"

	TitleDeviceDir    Title = "Device"
	TitleLanguageDir  Title = "Language"
	TitleViewModeDir  Title = "View Mode"
	TitleFloatingDir  Title = "Floating"
	TitleClipboardDir Title = "Clipboard"
	TitleMetricsDir   Title = "Metrics"

	TitleSelectedDevice    Title = "Selected Device"
	TitleSelectedLanguage  Title = "Selected Language"
	TitleSelectedViewMode  Title = "Selected View Mode"
	TitleSelectedFloating  Title = "Selected Floating"
	TitleSelectedClipboard Title = "Selected Clipboard"

	TitleOpenSettings  Title = "Open Settings"
	TitleStreamingText Title = "Streaming Text"
)
