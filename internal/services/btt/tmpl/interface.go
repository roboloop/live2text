package tmpl

//go:generate minimock -g -i Renderer -s _mock.go -o .

type Renderer interface {
	PrintStatus() string
	PrintMetric(template MetricTemplate, metric string, title string) string

	PrintSelectedDevice() string
	PrintSelectedLanguage() string
	PrintSelectedViewMode() string
	PrintSelectedFloating() string
	PrintSelectedClipboard() string

	SelectDevice(device string) string
	SelectLanguage(language string) string
	SelectViewMode(viewMode string) string
	SelectFloating(floatingState string) string
	SelectClipboard(clipboard string) string

	FloatingPage() string
	CloseSettings(
		cleanViewMode string,
		closeAction map[string]any,
		openCleanViewAction map[string]any,
		appUUID string,
		refreshAppPayload map[string]any,
	) string
	OpenSettings(action map[string]any) string
	Toggle() string
	ListenSocket(socketPath string) string
	AppPlaceholder() string
	CopyText() string
}
