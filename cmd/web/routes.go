package main

import (
	"github.com/bmizerany/pat"
	"github.com/zahnah/study-app/pkg/config"
	"github.com/zahnah/study-app/pkg/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	return mux
}
