package templates

import (
	"embed"
	"text/template"
)

//go:embed *.tmpl
var templates embed.FS

// Parse parses declared templates.
func Parse(t *template.Template) (*template.Template, error) {
	if t == nil {
		t = template.New("goapi-gen")
	}
	return t.ParseFS(templates, "*.tmpl")
}
