package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zahnah/study-app/internal/handlers"
	"net/http"
)

func routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoServe)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals", handlers.Repo.Generals)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/majors", handlers.Repo.Majors)

	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.PostAvailabilityJSON)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book-room", handlers.Repo.BookRoom)

	mux.Get("/user/login", handlers.Repo.Login)
	mux.Post("/user/login", handlers.Repo.PostLogin)
	mux.Get("/user/logout", handlers.Repo.Logout)

	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Post("/make-reservation", handlers.Repo.PostMakeReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Route("/admin", func(r chi.Router) {
		// temporary disable
		// r.Use(Auth)
		r.Get("/", handlers.Repo.AdminDashboard)
		r.Get("/reservations", handlers.Repo.AdminReservations)
		r.Get("/reservations/new", handlers.Repo.AdminReservationsNew)
		r.Get("/reservations/calendar", handlers.Repo.AdminReservationsCalendar)
		r.Get("/reservations/{src}/{id}", handlers.Repo.AdminReservation)
		r.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostReservation)
		r.Post("/reservations/{src}/{id}/processed", handlers.Repo.AdminProcessedReservation)
	})

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
