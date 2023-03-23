package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/zahnah/study-app/internal/config"
	"github.com/zahnah/study-app/internal/forms"
	"github.com/zahnah/study-app/internal/helpers"
	"github.com/zahnah/study-app/internal/models"
	"github.com/zahnah/study-app/internal/render"
	"github.com/zahnah/study-app/repository"
	"github.com/zahnah/study-app/repository/dbrepo"
	"net/http"
	"strconv"
	"time"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(app *config.AppConfig, db *sql.DB) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewPostgresRepo(db, app),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(writer http.ResponseWriter, r *http.Request) {
	render.Template(writer, *r, "home.page.gohtml", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(writer http.ResponseWriter, r *http.Request) {
	render.Template(writer, *r, "about.page.gohtml", &models.TemplateData{})
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

	sd := r.Form.Get("start")
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	ed := r.Form.Get("end")
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(writer, r, "/search-availability", http.StatusSeeOther)
		return
	} else {
		for _, i := range rooms {
			m.App.InfoLog.Println("ROOM", i.ID, i.RoomName)
		}
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)

	_ = render.Template(writer, *r, "choose-room.page.gohtml", &models.TemplateData{
		Data: data,
	})
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
		helpers.ServerError(writer, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(out)
}

func (m *Repository) MakeReservation(writer http.ResponseWriter, request *http.Request) {

	res, ok := m.App.Session.Get(request.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(writer, errors.New("cannot get reservation from session"))
		return
	}

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	room, err := m.DB.GetRoomById(res.RoomID)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}
	res.Room = room

	_ = render.Template(writer, *request, "make-reservation.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: map[string]interface{}{
			"reservation": res,
		},
		StringMap: map[string]string{
			"StartDate": sd,
			"EndDate":   ed,
		},
	})
}

func (m *Repository) PostMakeReservation(writer http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	// 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)
	data := make(map[string]interface{})
	data["reservation"] = reservation

	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(writer, *r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})
	} else {
		newID, err := m.DB.InsertReservation(reservation)
		if err != nil {
			helpers.ServerError(writer, err)
			return
		}

		restriction := models.RoomRestriction{
			RestrictionID: 1,
			ReservationID: newID,
			RoomID:        roomID,
			StartDate:     startDate,
			EndDate:       endDate,
		}

		_, err = m.DB.InsertRoomRestriction(restriction)
		if err != nil {
			helpers.ServerError(writer, err)
			return
		}

		m.App.Session.Put(r.Context(), "flash", "Data stored successfully")
		m.App.Session.Put(r.Context(), "reservation", reservation)
		http.Redirect(writer, r, "/reservation-summary", http.StatusSeeOther)
	}

}

func (m *Repository) ReservationSummary(writer http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		err := errors.New("can't get item from session")
		helpers.ServerError(writer, err)
		return
	} else {
		m.App.Session.Remove(r.Context(), "reservation")

		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(writer, *r, "reservation-summary.page.gohtml", &models.TemplateData{
			Data: data,
		})
	}
}

func (m *Repository) ChooseRoom(writer http.ResponseWriter, request *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		helpers.ServerError(writer, err)
		return
	}

	res, ok := m.App.Session.Get(request.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(writer, err)
		return
	}

	res.RoomID = roomID
	m.App.Session.Put(request.Context(), "reservation", res)

	http.Redirect(writer, request, "/make-reservation", http.StatusSeeOther)
}
