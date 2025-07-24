package trigger

func NewFloatingMenu(title Title) Trigger {
	return newTrigger().
		initFloatingMenu(title, triggerFloatingMenu, typeFloatingMenu).
		// TODO: btt bug, need to manually update the trigger to make it visible
		addExtra(map[string]any{
			"BTTMenuConfig": map[string]any{
				"BTTMenuVisibility": 1,

				"BTTMenuCategoryMenuVisibility": 1,
				"BTTMenuCategorySize":           0,
				"BTTMenuSizingBehavior":         3,

				"BTTMenuPositioningType":                     1,
				"BTTMenuAnchorMenu":                          4,
				"BTTMenuPositionRelativeTo":                  23,
				"BTTMenuAnchorRelation":                      5,
				"BTTMenuOffsetX":                             0.0,
				"BTTMenuOffsetXUnit":                         1,
				"BTTMenuOffsetY":                             50.0,
				"BTTMenuOffsetYUnit":                         1,
				"BTTMenuCategoryPosition":                    0,
				"BTTMenuOnlyUpdatePositionOnExplicitRequest": 1,

				"BTTMenuItemPaddingTop":    25,
				"BTTMenuItemPaddingBottom": 0,
				"BTTMenuItemPaddingLeft":   0,
				"BTTMenuItemPaddingRight":  0,
				"BTTMenuVerticalSpacing":   0,
				"BTTMenuHorizontalSpacing": 0,

				"BTTMenuOpacityActive":   1.0,
				"BTTMenuOpacityInactive": 0.9,

				"BTTMenuItemBlurredBackground": 0,
				"BTTMenuAppearanceStyle":       0,

				"BTTMenuItemBorderWidth": 0,

				"BTTMenuSelectedTab":        0,
				"BTTMenuAvailability":       -1,
				"BTTMenuItemSelectedTab":    0,
				"BTTMenuUseStyleForSubmenu": 1,
			},
		})
}

func NewWebView(title Title, content string) Trigger {
	return newTrigger().
		initFloatingMenu(title, triggerWebView, typeFloatingMenu).
		addExtra(map[string]any{
			"BTTMenuConfig": map[string]any{
				"BTTMenuVisibility": 1,

				"BTTMenuItemMinWidth":  100,
				"BTTMenuItemMaxWidth":  550,
				"BTTMenuItemMinHeight": 50,
				"BTTMenuItemMaxHeight": 50,

				"BTTMenuItemVisibleWhileActive":   1,
				"BTTMenuItemVisibleWhileInactive": 1,

				"BTTMenuSizingBehavior": 0,

				"BTTMenuItemBlurredBackground": 0,
				"BTTMenuAppearanceStyle":       0,
				"BTTMenuItemUserAgent":         "BTT-CLIENT",

				"BTTMenuItemBorderWidth":  0,
				"BTTMenuItemCornerRadius": 10,

				"BTTMenuUseStyleForSubmenu": 1,
				"BTTMenuItemSelectedTab":    0,

				"BTTMenuItemText": content,
			},
		})
}

func (t Trigger) initFloatingMenu(title Title, triggerID id, triggerType bttType) Trigger {
	return t.addExtra(map[string]any{
		"BTTMenuName":     title.String(),
		"BTTTriggerType":  triggerID.float(),
		"BTTTriggerClass": triggerType.String(),
		"BTTMenuConfig": map[string]any{
			"BTTMenuItemShowIdentifierAsTooltip": 0,
			"BTTMenuItemPositioningMode":         0,

			"BTTMenuAlwaysUseLightMode": 1,

			"BTTMenuElementIdentifier": title.String(),
		},
	})
}
