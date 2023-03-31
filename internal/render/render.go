package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/justinas/nosurf"
	"github.com/zahnah/study-app/internal/config"
	"github.com/zahnah/study-app/internal/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var functions = template.FuncMap{
	"humanDate": HumanDate,
}

var app *config.AppConfig
var pathToTemplates = "./templates"

func NewRenderer(a *config.AppConfig) {
	app = a
}

func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")

	td.CSRFToken = nosurf.Token(r)
	td.IsAuthenticated = app.Session.Exists(r.Context(), "user_id")

	return td
}

func Template(writer http.ResponseWriter, r http.Request, tmpl string, td *models.TemplateData) error {

	tc := app.TemplateCache
	if !app.UseCache {
		tc, _ = CreateTemplateCache()
	}

	parseTemplate, ok := tc[tmpl]
	if !ok {
		log.Println("Couldn't get template for the page")
		return errors.New("couldn't get template for the page")
	}

	buf := new(bytes.Buffer)

	_ = parseTemplate.Execute(buf, AddDefaultData(td, &r))

	_, err := buf.WriteTo(writer)

	if err != nil {
		fmt.Println("error parsing template:", err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
