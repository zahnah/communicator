package handlers

import (
	"github.com/zahnah/study-app/pkg/render"
	"net/http"
)

// Home is the home page handler
func Home(writer http.ResponseWriter, r *http.Request) {
	render.Template(writer, "home.page.gohtml")
}

// About is the about page handler
func About(writer http.ResponseWriter, r *http.Request) {
	render.Template(writer, "about.page.gohtml")
}
