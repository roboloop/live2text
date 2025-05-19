package payload

func (p Payload) AddIcon(sfSymbol string, height int, onlyIcon bool) Payload {
	p.AddMap(map[string]any{
		"BTTTriggerConfig": map[string]any{
			"BTTTouchBarItemSFSymbolDefaultIcon": sfSymbol,
			"BTTTouchBarItemSFSymbolWeight":      0,
			"BTTTouchBarItemIconType":            2,
			"BTTTouchBarItemIconHeight":          height,
			"BTTTouchBarItemPadding":             -10,
		},
	})

	if onlyIcon {
		p.AddMap(map[string]any{
			"BTTTriggerConfig": map[string]any{
				"BTTTouchBarButtonColor":  "0.0, 0.0, 0.0, 255.0",
				"BTTTouchBarOnlyShowIcon": true,
			},
		})
	}

	return p
}
