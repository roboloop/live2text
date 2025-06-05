package tmpl

import (
	"bytes"
	"fmt"
	"maps"
	"text/template"

	_ "embed"
)

//go:embed scripts.shell.gotmpl
var scripts []byte

//go:embed pages.html.gotmpl
var pages []byte

type Renderer struct {
	tmpl    *template.Template
	appName string
	debug   bool
}

func NewRenderer(appName string, debug bool) *Renderer {
	tmpl := template.New("templates")
	tmpl = template.Must(tmpl.Parse(string(scripts)))
	tmpl = template.Must(tmpl.Parse(string(pages)))

	return &Renderer{tmpl: tmpl, appName: appName, debug: debug}
}

func (r *Renderer) Render(name string, data map[string]any) (string, error) {
	cloned := maps.Clone(data)
	cloned["Debug"] = r.debug
	cloned["Name"] = name
	cloned["AppName"] = r.appName

	var buf bytes.Buffer
	if err := r.tmpl.ExecuteTemplate(&buf, name, cloned); err != nil {
		return "", fmt.Errorf("cannot execute template: %w", err)
	}

	return buf.String(), nil
}
