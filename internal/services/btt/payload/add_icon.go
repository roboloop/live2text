package payload

const DefaultHeight = 22

const (
	IconMicrophone = "microphone"
	IconCharacter  = "character"
	IconMacwindow  = "macwindow"
	IconChartBar   = "chart.bar"

	IconWaveform                   = "waveform"
	IconSquareAndArrowUp           = "square.and.arrow.up"
	IconSquareAndArrowUpBadgeClock = "square.and.arrow.up.badge.clock"
	IconFlame                      = "flame"
	IconAppConnectedToAppBelowFill = "app.connected.to.app.below.fill"
	IconListBullet                 = "list.bullet"
	IconRectangleStack             = "rectangle.stack"
)

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
