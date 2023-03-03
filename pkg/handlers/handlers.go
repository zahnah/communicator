package handlers

import (
	"github.com/zahnah/study-app/pkg/config"
	"github.com/zahnah/study-app/pkg/models"
	"github.com/zahnah/study-app/pkg/render"
	"net/http"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(writer http.ResponseWriter, r *http.Request) {
	render.Template(writer, "home.page.gohtml", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(writer http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	render.Template(writer, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}
