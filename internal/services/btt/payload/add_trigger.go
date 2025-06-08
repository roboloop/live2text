package payload

type TriggerID int

const (
	TriggerDirectory      TriggerID = 630
	TriggerTouchBarButton TriggerID = 629
	TriggerNamed          TriggerID = 643
	TriggerShellScript    TriggerID = 642
	TriggerFloatingMenu   TriggerID = 767
	TriggerWebView        TriggerID = 778
)

type TriggerType string

const (
	TouchBar      TriggerType = "BTTTriggerTypeTouchBar"
	OtherTriggers TriggerType = "BTTTriggerTypeOtherTriggers"
	FloatingMenu  TriggerType = "BTTTriggerTypeFloatingMenu"
)

type ActionType int

const (
	ActionTypeNone                          = 0
	ActionTypeEmptyPlaceholder   ActionType = 366
	ActionTypeExecuteScript      ActionType = 206
	ActionTypeCloseGroup         ActionType = 191
	ActionTypeOpenGroup          ActionType = 205
	ActionTypeOpenFloatingMenu   ActionType = 386
	ActionTypeCloseFloatingMenu  ActionType = 387
	ActionTypeToggleFloatingMenu ActionType = 388
)

func (p Payload) AddTrigger(
	name string,
	triggerID TriggerID,
	triggerType TriggerType,
	actionType ActionType,
	hidden bool,
) Payload {
	p.AddMap(map[string]any{
		"BTTTouchBarButtonName": name,
		"BTTWidgetName":         name,
		"BTTTriggerName":        name,
		"BTTMenuName":           name,
		"BTTTriggerType":        triggerID,
		"BTTTriggerClass":       triggerType,
		"BTTGroupName":          LabelName,
		"BTTNotes":              LabelName,
		"BTTTriggerConfig": map[string]any{
			"BTTKeepGroupOpenWhileSwitchingApps": true,
		},
	})

	if actionType != ActionTypeNone {
		p.AddMap(map[string]any{
			"BTTPredefinedActionType": actionType,
		})
	}

	if triggerType == OtherTriggers {
		p.AddMap(map[string]any{
			"BTTGestureNotes": LabelName,
		})
	}

	if hidden {
		p.AddMap(map[string]any{
			"BTTTriggerConfig": map[string]any{
				"BTTTouchBarButtonWidth":         0,
				"BTTTouchBarButtonUseFixedWidth": 1,
			},
		})
	}

	return p
}
