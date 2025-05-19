package tmpl

import (
	"bytes"
	_ "embed"
	"fmt"
	"maps"
	"text/template"
)

//go:embed scripts.shell.gotmpl
var scripts []byte

type Renderer struct {
	tmpl    *template.Template
	appName string
	debug   bool
}

func NewRenderer(appName string, debug bool) *Renderer {
	tmpl := template.Must(template.New("scripts").Parse(string(scripts)))
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
