package trigger

const defaultHeight = 22

type Icon string

func (i Icon) String() string {
	return string(i)
}

const (
	IconMicrophone Icon = "microphone"
	IconCharacter  Icon = "character"
	IconMacwindow  Icon = "macwindow"
	IconChartBar   Icon = "chart.bar"

	IconXmarkCircleFill            Icon = "xmark.circle.fill"
	IconWaveform                   Icon = "waveform"
	IconSquareAndArrowUp           Icon = "square.and.arrow.up"
	IconSquareAndArrowUpBadgeClock Icon = "square.and.arrow.up.badge.clock"
	IconFlame                      Icon = "flame"
	IconAppConnectedToAppBelowFill Icon = "app.connected.to.app.below.fill"
	IconListBullet                 Icon = "list.bullet"
	IconRectangleStack             Icon = "rectangle.stack"
)

func (t Trigger) addIcon(sfSymbol Icon, height int, onlyIcon bool) Trigger {
	t.addExtra(map[string]any{
		"BTTTriggerConfig": map[string]any{
			"BTTTouchBarItemSFSymbolDefaultIcon": sfSymbol.String(),
			"BTTTouchBarItemSFSymbolWeight":      0,
			"BTTTouchBarItemIconType":            2,
			"BTTTouchBarItemIconHeight":          float64(height),
			"BTTTouchBarItemPadding":             -10,
		},
	})

	if onlyIcon {
		t.addExtra(map[string]any{
			"BTTTriggerConfig": map[string]any{
				"BTTTouchBarButtonColor":  "0.0, 0.0, 0.0, 255.0",
				"BTTTouchBarOnlyShowIcon": true,
			},
		})
	}

	return t
}
