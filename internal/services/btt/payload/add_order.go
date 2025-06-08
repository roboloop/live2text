package payload

type Order int

const (
	MainOrderSettings Order = 20
	MainOrderViewMode Order = 21
	MainOrderSubs     Order = 22
)

const (
	NamedOrderSettings Order = iota
)

const (
	SettingsOrderCloseGroup Order = iota
	SettingsOrderStatus
	SettingsOrderDevice
	SettingsOrderLanguage
	SettingsOrderViewMode
	SettingsOrderFloating
	SettingsOrderMetrics
)

const (
	DeviceOrderCloseGroup Order = iota
	DeviceOrderSelectedDevice
)

const (
	LanguageOrderCloseGroup Order = iota
	LanguageOrderSelectedLanguage
)

const (
	ViewModeOrderCloseGroup Order = iota
	ViewModeOrderSelectedViewMode
)

const (
	FloatingOrderCloseGroup Order = iota
	FloatingOrderSelectedState
)

const (
	MetricsOrderCloseGroup Order = iota
	MetricsOrderReadAudio
	MetricsOrderSentAudio
	MetricsOrderSendAudioMs
	MetricsOrderBurntAudio
	MetricsOrderConnections
	MetricsOrderTasks
	MetricsOrderSockets
)

func (p Payload) AddOrder(order Order) Payload {
	p.AddMap(map[string]any{
		"BTTOrder": order,
	})

	return p
}
