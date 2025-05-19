package payload

type Order int

const (
	MainOrderSettings Order = 20
	MainOrderSubs     Order = 21
)

const (
	NamedOrderSettings Order = iota
)

const (
	SettingsOrderCloseGroup Order = iota
	SettingsOrderStatus
	SettingsOrderDevice
	SettingsOrderLanguage
)

const (
	DeviceOrderCloseGroup Order = iota
	DeviceOrderSelectedDevice
)

const (
	LanguageOrderCloseGroup Order = iota
	LanguageOrderSelectedLanguage
)

func (p Payload) AddOrder(order Order) Payload {
	p.AddMap(map[string]any{
		"BTTOrder": order,
	})

	return p
}
