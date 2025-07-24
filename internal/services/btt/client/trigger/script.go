package trigger

const (
	settingsInterval = 15
	metricsInterval  = 5
)

func (t Trigger) AddTapScript(script string) Trigger {
	return t.addTapScript(script)
}

func (t Trigger) AddLongTapTrigger(triggerTitle Title) Trigger {
	return t.addLongTapScript(triggerTitle)
}

func (t Trigger) HasTapScript() bool {
	rawActions, ok1 := t["BTTAdditionalActions"]
	if !ok1 {
		return false
	}

	actions, ok := rawActions.([]any)
	if !ok || len(actions) == 0 {
		return false
	}

	value, ok := actions[0].(map[string]any)["BTTPredefinedActionType"]
	if !ok {
		return false
	}

	floatVal, ok := value.(float64)

	return ok && floatVal == actionTypeExecuteScript.Float()
}

func (t Trigger) AddCycledScript(script string, interval float64) Trigger {
	return t.addCycledScript(script, interval)
}

func (t Trigger) addCycledScript(script string, interval float64) Trigger {
	trigger := map[string]any{
		"BTTShellScriptWidgetGestureConfig": "/bin/bash:::-c:::-:::",
		"BTTTriggerConfig": map[string]any{
			"BTTTouchBarAppleScriptStringRunOnInit": 1,
			"BTTTouchBarShellScriptString":          script,
			"BTTTouchBarScriptUpdateInterval":       interval,
		},
	}
	if interval == 0.0 {
		trigger["BTTTriggerConfig"].(map[string]any)["BTTTouchBarAppleScriptStringRunOnInit"] = 0
	}

	return t.addExtra(trigger)
}

func (t Trigger) addTapScript(script string) Trigger {
	trigger := map[string]any{
		"BTTAdditionalActions": []any{
			map[string]any{
				"BTTPredefinedActionType":  actionTypeExecuteScript.Float(),
				"BTTShellTaskActionScript": script,
				"BTTShellTaskActionConfig": "/bin/bash:::-c:::-:::",
			},
		},
	}

	return t.addExtra(trigger)
}

func (t Trigger) addLongTapScript(triggerTitle Title) Trigger {
	return t.addExtra(map[string]any{
		"BTTTriggerConfig": map[string]any{
			"BTTTouchBarLongPressActionName": triggerTitle.String(),
		},
	})
}
