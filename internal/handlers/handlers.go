package handlers

import (
	"encoding/json"
	"github.com/zahnah/study-app/internal/config"
	"github.com/zahnah/study-app/internal/models"
	"github.com/zahnah/study-app/internal/render"
	"log"
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
	render.Template(writer, *r, "home.page.gohtml", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(writer http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.Template(writer, *r, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Generals(writer http.ResponseWriter, r *http.Request) {
	render.Template(writer, *r, "generals.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Contact(writer http.ResponseWriter, r *http.Request) {
	render.Template(writer, *r, "contact.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Majors(writer http.ResponseWriter, r *http.Request) {
	render.Template(writer, *r, "majors.page.gohtml", &models.TemplateData{})
}

func (m *Repository) SearchAvailability(writer http.ResponseWriter, r *http.Request) {
	render.Template(writer, *r, "search-availability.page.gohtml", &models.TemplateData{})
}

func (m *Repository) PostAvailability(writer http.ResponseWriter, r *http.Request) {
	_, _ = writer.Write([]byte("Post availability"))
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	_, _ = writer.Write([]byte(start + " - " + end))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) PostAvailabilityJSON(writer http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")

	if err != nil {
		log.Println(err)
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(out)
}
