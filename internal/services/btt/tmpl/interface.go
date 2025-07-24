package tmpl

//go:generate minimock -g -i Renderer -s _mock.go -o .

type Renderer interface {
	PrintStatus() string
	PrintMetric(template MetricTemplate, metric string, title string) string

	PrintSelectedDevice() string
	PrintSelectedLanguage() string
	PrintSelectedViewMode() string
	PrintSelectedFloatingState() string

	SelectDevice(device string) string
	SelectLanguage(language string) string
	SelectViewMode(viewMode string) string
	SelectFloatingState(floatingState string) string

	FloatingPage() string
	CloseSettings(cleanViewMode string, closeAction map[string]any, openCleanViewAction map[string]any) string
	OpenSettings(action map[string]any) string
	Toggle() string
	ListenSocket(socketPath string) string
	AppPlaceholder() string
}
