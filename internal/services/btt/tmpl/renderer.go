package tmpl

import (
	"text/template"

	_ "embed"
)

//go:embed templates/common.shell.gotmpl
var common []byte

//go:embed templates/scripts.shell.gotmpl
var scripts []byte

//go:embed templates/metrics.shell.gotmpl
var metrics []byte

//go:embed templates/pages.html.gotmpl
var pages []byte

type renderer struct {
	tmpl       *template.Template
	appName    string
	appAddress string
	bttAddress string
	debug      bool
}

func NewRenderer(appName string, appAddress, bttAddress string, debug bool) Renderer {
	tmpl := template.New("templates")
	tmpl = template.Must(tmpl.Parse(string(common)))
	tmpl = template.Must(tmpl.Parse(string(scripts)))
	tmpl = template.Must(tmpl.Parse(string(metrics)))
	tmpl = template.Must(tmpl.Parse(string(pages)))

	return &renderer{
		tmpl:       tmpl,
		appAddress: appAddress,
		bttAddress: bttAddress,
		appName:    appName,
		debug:      debug,
	}
}
