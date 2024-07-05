package app

import (
	"html/template"
	"io"
	"time"

	"github.com/labstack/echo/v4"
)

type Template struct {
	app       *App
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	r, err := t.templates.Clone()
	if err != nil {
		return err
	}

	r.Funcs(template.FuncMap{})

	return r.ExecuteTemplate(w, name, data)
}

func echoFunc(key string, _ ...interface{}) string {
	return key
}

func (a *App) viewTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"Version":   func() *Version { return &a.Version },
		"AppConfig": func() *Config { return &a.Config },
		"LocalTime": func(t time.Time) time.Time { return t.UTC() },
		"LocalDate": func(t time.Time) string { return t.UTC().Format("2006-01-02 15:04") },

		"RouteFor": func(name string, params ...interface{}) string {
			rev := a.echo.Reverse(name, params...)
			if rev == "" {
				return "/invalid/route/#" + name
			}

			return rev
		},
	}
}
