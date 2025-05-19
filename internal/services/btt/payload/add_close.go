package payload

func (p Payload) AddClose(groupName string) Payload {
	p.AddTrigger("Close Group", TriggerTouchBarButton, TouchBar, ActionTypeOpenGroup, false)
	p.AddIcon("xmark.circle.fill", 25, true)
	p.AddMap(map[string]any{
		"BTTOpenGroupWithName": groupName,
	})

	return p
}
