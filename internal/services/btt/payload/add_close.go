package payload

func (p Payload) AddClose(groupName string) Payload {
	if groupName != "" {
		p.AddTrigger("Close Group", TriggerTouchBarButton, TouchBar, ActionTypeOpenGroup, false)
		p.AddMap(map[string]any{
			"BTTOpenGroupWithName": groupName,
		})
	} else {
		p.AddTrigger("Close Group", TriggerTouchBarButton, TouchBar, ActionTypeCloseGroup, false)
	}

	p.AddIcon("xmark.circle.fill", 25, true)

	return p
}
