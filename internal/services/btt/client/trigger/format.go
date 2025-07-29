package trigger

func (t Trigger) AddReadableFormat() Trigger {
	return t.addExtra(map[string]any{
		"BTTTriggerConfig": map[string]any{
			"BTTTouchBarButtonFontSize":        12,
			"BTTTouchBarButtonColor":           "0.0, 0.0, 0.0, 255.0",
			"BTTTouchBarButtonTextAlignment":   0,
			"BTTShellScriptDontTrimWhitepsace": 1,
			"BTTShellScriptDontTrimWhitespace": 1,
		},
	})
}

func (t Trigger) AddEnabled() Trigger {
	return t.addExtra(map[string]any{
		"BTTEnabled2": 1,
	})
}

func (t Trigger) AddDisabled() Trigger {
	return t.addExtra(map[string]any{
		"BTTEnabled2": 0,
	})
}
