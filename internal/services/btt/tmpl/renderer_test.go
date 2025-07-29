package tmpl_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"live2text/internal/services/btt/tmpl"
)

func TestPrintStatus(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.PrintStatus()
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestPrintMetric(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)

	result := renderer.PrintMetric(tmpl.MetricTemplateRaw, "app_foo", "bar")
	require.NotEmpty(t, strings.TrimSpace(result))

	result = renderer.PrintMetric(tmpl.MetricTemplateSize, "app_bar", "bar")
	require.NotEmpty(t, strings.TrimSpace(result))

	result = renderer.PrintMetric(tmpl.MetricTemplateDuration, "app_baz", "baz")
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestPrintSelectedDevice(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.PrintSelectedDevice()
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestPrintSelectedLanguage(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.PrintSelectedLanguage()
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestPrintSelectedViewMode(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.PrintSelectedViewMode()
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestPrintSelectedFloating(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.PrintSelectedFloating()
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestPrintSelectedClipboard(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.PrintSelectedClipboard()
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestSelectDevice(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.SelectDevice("foo")
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestSelectLanguage(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.SelectLanguage("foo")
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestSelectViewMode(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.SelectViewMode("foo")
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestSelectFloating(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.SelectFloating("foo")
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestSelectClipboard(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.SelectClipboard("foo")
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestFloatingPage(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.FloatingPage()
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestOpenSettings(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.OpenSettings(map[string]any{"foo": "bar"})
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestCloseSettings(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.CloseSettings("Clean", map[string]any{"foo": "bar"}, map[string]any{"key": "value"}, "", nil)
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestToggle(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.Toggle()
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestListenSocket(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.ListenSocket("/path/to/file.sock")
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestAppPlaceholder(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.AppPlaceholder()
	require.NotEmpty(t, strings.TrimSpace(result))
}

func TestCopyText(t *testing.T) {
	t.Parallel()

	renderer := setupRenderer(t)
	result := renderer.CopyText()
	require.NotEmpty(t, strings.TrimSpace(result))
}

func setupRenderer(t *testing.T) tmpl.Renderer {
	t.Helper()

	return tmpl.NewRenderer("FooApp", "127.0.0.1:1010", "192.168.1.1:2020", false)
}
