package trigger

type UUID string

func (uuid UUID) String() string {
	return string(uuid)
}

func (t Trigger) AddUUID(uuid UUID) Trigger {
	return t.addExtra(map[string]any{
		"BTTUUID": uuid.String(),
	})
}

func (t Trigger) UUID() UUID {
	if v, ok := t["BTTUUID"].(string); ok {
		return UUID(v)
	}

	return ""
}
