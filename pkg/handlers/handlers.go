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
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)
	render.Template(writer, "home.page.gohtml", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(writer http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.Template(writer, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Generals(writer http.ResponseWriter, r *http.Request) {
	render.Template(writer, "generals.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Contact(writer http.ResponseWriter, r *http.Request) {
	render.Template(writer, "contact.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Majors(writer http.ResponseWriter, r *http.Request) {
	render.Template(writer, "majors.page.gohtml", &models.TemplateData{})
}

func (m *Repository) SearchAvailability(writer http.ResponseWriter, r *http.Request) {
	render.Template(writer, "search-availability.page.gohtml", &models.TemplateData{})
}
