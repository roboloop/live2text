package tmpl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"maps"
	"net/url"
	"strings"
)

type MetricTemplate string

const (
	MetricTemplateSize     MetricTemplate = "print_size_metric"
	MetricTemplateDuration MetricTemplate = "print_duration_metric"
	MetricTemplateRaw      MetricTemplate = "print_raw_metric"
)

func (r *renderer) render(name string, data map[string]any) string {
	cloned := maps.Clone(data)
	cloned["Debug"] = r.debug
	cloned["DefinitionName"] = name
	cloned["AppName"] = r.appName
	cloned["AppAddress"] = r.appAddress
	cloned["BttAddress"] = r.bttAddress

	var buf bytes.Buffer
	if err := r.tmpl.ExecuteTemplate(&buf, name, cloned); err != nil {
		panic(fmt.Errorf("cannot execute template: %w", err))
	}

	return buf.String()
}

func (r *renderer) PrintStatus() string {
	return r.render("print_status", map[string]any{})
}

func (r *renderer) PrintMetric(template MetricTemplate, metric string, title string) string {
	return r.render(string(template), map[string]any{"Metric": metric, "Title": title})
}

func (r *renderer) PrintSelectedDevice() string {
	return r.render("print_selected_device", map[string]any{})
}

func (r *renderer) PrintSelectedLanguage() string {
	return r.render("print_selected_language", map[string]any{})
}

func (r *renderer) PrintSelectedViewMode() string {
	return r.render("print_selected_view_mode", map[string]any{})
}

func (r *renderer) PrintSelectedFloating() string {
	return r.render("print_selected_floating", map[string]any{})
}

func (r *renderer) PrintSelectedClipboard() string {
	return r.render("print_selected_clipboard", map[string]any{})
}

func (r *renderer) SelectDevice(device string) string {
	return r.render("select_device", map[string]any{"Device": device})
}

func (r *renderer) SelectLanguage(language string) string {
	return r.render("select_language", map[string]any{"Language": language})
}

func (r *renderer) SelectViewMode(viewMode string) string {
	return r.render("select_view_mode", map[string]any{"ViewMode": viewMode})
}

func (r *renderer) SelectFloating(floating string) string {
	return r.render("select_floating", map[string]any{"Floating": floating})
}

func (r *renderer) SelectClipboard(clipboard string) string {
	return r.render("select_clipboard", map[string]any{"Clipboard": clipboard})
}

func (r *renderer) FloatingPage() string {
	return r.render("floating_page", map[string]any{})
}

func (r *renderer) OpenSettings(action map[string]any) string {
	query, err := encodeForShell(action)
	if err != nil {
		return ""
	}

	return r.render("open_settings", map[string]any{
		"Query": query,
	})
}

func (r *renderer) CloseSettings(
	cleanViewMode string,
	closeAction map[string]any,
	openCleanViewAction map[string]any,
	appUUID string,
	refreshAppPayload map[string]any,
) string {
	closeQuery, err := encodeForShell(closeAction)
	if err != nil {
		return ""
	}
	openCleanViewQuery, err := encodeForShell(openCleanViewAction)
	if err != nil {
		return ""
	}
	refreshAppQuery, err := encodeForShell(refreshAppPayload)
	if err != nil {
		return ""
	}

	return r.render("close_settings", map[string]any{
		"CloseGroupQuery":    closeQuery,
		"OpenCleanViewQuery": openCleanViewQuery,
		"CleanViewMode":      cleanViewMode,

		"AppUUID":         appUUID,
		"RefreshAppQuery": refreshAppQuery,
	})
}

func (r *renderer) Toggle() string {
	return r.render("toggle", map[string]any{})
}

func (r *renderer) ListenSocket(socketPath string) string {
	return r.render("listen_socket", map[string]any{"SocketPath": socketPath})
}

func (r *renderer) AppPlaceholder() string {
	return r.render("app_placeholder", map[string]any{})
}

func (r *renderer) CopyText() string {
	return r.render("copy_text", map[string]any{})
}

func encodeForShell(jsonPayload map[string]any) (string, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(jsonPayload); err != nil {
		return "", fmt.Errorf("cannot encode json payload: %w", err)
	}

	query := url.Values{}
	query.Set("json", buf.String())
	encoded := strings.ReplaceAll(query.Encode(), "+", "%20")

	return encoded, nil
}
