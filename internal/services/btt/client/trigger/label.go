package trigger

func (t Trigger) AddLabel(label string) Trigger {
	t.addExtra(map[string]any{
		"BTTGroupName": label,
		"BTTNotes":     label,
	})

	if t["BTTTriggerClass"] == typeOtherTriggers {
		t.addExtra(map[string]any{
			"BTTGestureNotes": label,
		})
	}

	return t
}
