package trigger

const (
	OrderSettings = 20
)

func (t Trigger) AddOrder(order int) Trigger {
	return t.addExtra(map[string]any{
		"BTTOrder": float64(order),
	})
}

func (t Trigger) AddOrderAfter(after Trigger) Trigger {
	if after == nil {
		t["BTTOrder"] = float64(0)
		return t
	}

	order, _ := after["BTTOrder"].(float64)
	t["BTTOrder"] = float64(order + 1)

	return t
}
