package payload

func (p Payload) AddParentUUID(uuid string) Payload {
	p.AddMap(map[string]any{
		"BTTTriggerParentUUID": uuid,
	})

	return p
}
