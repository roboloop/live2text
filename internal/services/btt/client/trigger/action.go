package trigger

func NewCloseDirAction() Trigger {
	return newTrigger().addCloseDirAction()
}

func NewOpenDirAction(dirName Title) Trigger {
	return newTrigger().addOpenDirAction(dirName)
}

func NewOpenFloatingAction(title Title) Trigger {
	return newTrigger().
		addAction(actionTypeOpenFloatingMenu).
		addExtra(map[string]any{
			"BTTAdditionalActionData": map[string]any{
				"BTTMenuActionMenuID": title.String(),
			},
		})
}

func NewCloseFloatingAction(title Title) Trigger {
	return newTrigger().
		addAction(actionTypeCloseFloatingMenu).
		addExtra(map[string]any{
			"BTTAdditionalActionData": map[string]any{
				"BTTMenuActionMenuID": title.String(),
			},
		})
}

func (t Trigger) addCloseDirAction() Trigger {
	return t.addAction(actionTypeCloseGroup)
}

func (t Trigger) addOpenDirAction(dirName Title) Trigger {
	return t.addAction(actionTypeOpenGroup).addExtra(map[string]any{
		"BTTOpenGroupWithName": dirName.String(),
	})
}

func (t Trigger) addAction(actionType actionType) Trigger {
	return t.addExtra(map[string]any{
		"BTTPredefinedActionType": actionType.Float(),
	})
}
