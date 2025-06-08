package payload

type ShellType int

const (
	ShellTypeNone ShellType = iota
	ShellTypeEmbed
	ShellTypeAdditional
)

const (
	IntervalDefault = 15
	IntervalMetrics = 5
	IntervalNone    = 0
)

func (p Payload) AddShell(shell string, interval float64, shellType ShellType) Payload {
	trigger := map[string]any{
		"BTTShellScriptWidgetGestureConfig": "/bin/bash:::-c:::-:::",
		"BTTTriggerConfig": map[string]any{
			"BTTTouchBarAppleScriptStringRunOnInit": 1,
			"BTTTouchBarShellScriptString":          shell,
		},
	}

	embed := map[string]any{
		"BTTPredefinedActionType":  ActionTypeExecuteScript,
		"BTTShellTaskActionScript": shell,
		"BTTShellTaskActionConfig": "/bin/bash:::-c:::-:::",
	}

	additional := map[string]any{
		"BTTAdditionalActions": []map[string]any{embed},
	}

	if interval > 0.0 {
		trigger["BTTTriggerConfig"].(map[string]any)["BTTTouchBarScriptUpdateInterval"] = interval
	}

	switch shellType {
	case ShellTypeNone:
		p.AddMap(trigger)
	case ShellTypeEmbed:
		p.AddMap(embed)
	case ShellTypeAdditional:
		p.AddMap(additional)
	}

	return p
}
