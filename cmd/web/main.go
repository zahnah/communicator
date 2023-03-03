package main

import (
	"fmt"
	"github.com/zahnah/study-app/pkg/config"
	"github.com/zahnah/study-app/pkg/handlers"
	"github.com/zahnah/study-app/pkg/render"
	"log"
	"net/http"
)

const portNumber = ":8080"

func main() {

	var app config.AppConfig
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatalln("Can't create a template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println(fmt.Sprintf("Starting application on port: %s", portNumber))
	_ = http.ListenAndServe(portNumber, nil)
}
