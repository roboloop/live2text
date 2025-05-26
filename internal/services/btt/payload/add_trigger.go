package payload

type TriggerID int

const (
	TriggerDirectory      TriggerID = 630
	TriggerTouchBarButton TriggerID = 629
	TriggerNamed          TriggerID = 643
	TriggerShellScript    TriggerID = 642
)

type TriggerType string

const (
	TouchBar      TriggerType = "BTTTriggerTypeTouchBar"
	OtherTriggers TriggerType = "BTTTriggerTypeOtherTriggers"
)

type ActionType int

const (
	ActionTypeEmptyPlaceholder ActionType = 366
	ActionTypeExecuteScript    ActionType = 206
	ActionTypeCloseGroup       ActionType = 191
	ActionTypeOpenGroup        ActionType = 205
)

func (p Payload) AddTrigger(
	name string,
	triggerID TriggerID,
	triggerType TriggerType,
	actionType ActionType,
	hidden bool,
) Payload {
	p.AddMap(map[string]any{
		"BTTTouchBarButtonName":   name,
		"BTTWidgetName":           name,
		"BTTTriggerName":          name,
		"BTTTriggerType":          triggerID,
		"BTTTriggerClass":         triggerType,
		"BTTPredefinedActionType": actionType,
		"BTTGroupName":            noteName,
		"BTTNotes":                noteName,
		"BTTTriggerConfig": map[string]any{
			"BTTKeepGroupOpenWhileSwitchingApps": true,
		},
	})

	if triggerType == OtherTriggers {
		p.AddMap(map[string]any{
			"BTTGestureNotes": noteName,
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
