package trigger

import "fmt"

type Trigger map[string]any

func NewTrigger() Trigger {
	return newTrigger()
}

func newTrigger() Trigger {
	return make(Trigger)
}

func (t Trigger) init(
	title Title,
	triggerID id,
	triggerType bttType,
) Trigger {
	return t.addExtra(map[string]any{
		"BTTTouchBarButtonName": title.String(),
		"BTTWidgetName":         title.String(),
		"BTTTriggerName":        title.String(),
		"BTTMenuName":           title.String(),
		"BTTTriggerType":        triggerID.float(),
		"BTTTriggerClass":       triggerType.String(),
		"BTTTriggerConfig": map[string]any{
			"BTTKeepGroupOpenWhileSwitchingApps": true,
		},
	})
}

func (t Trigger) ErrorContext() string {
	if _, ok := t["BTTTriggerName"]; ok {
		return fmt.Sprintf("trigger %s", t.Title())
	}

	if _, ok := t["BTTPredefinedActionType"]; ok {
		return fmt.Sprintf("action %s", t.actionType())
	}

	return "unknown trigger"
}

func (t Trigger) Title() Title {
	if v, ok := t["BTTTriggerName"].(string); ok {
		return Title(v)
	}

	return ""
}

func (t Trigger) actionType() actionType {
	if v, ok := t["BTTPredefinedActionType"].(float64); ok {
		return actionType(v)
	}

	return actionType(0.0)
}
