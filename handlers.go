package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Home is the home page handler
func Home(writer http.ResponseWriter, r *http.Request) {
	renderTemplate(writer, "home_page.html")
}

// About is the about page handler
func About(writer http.ResponseWriter, r *http.Request) {
	renderTemplate(writer, "about_page.html")
}

func renderTemplate(writer http.ResponseWriter, tmpl string) {
	parseTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	err := parseTemplate.Execute(writer, nil)
	if err != nil {
		fmt.Println("error parsing template:", err)
		return
	}
}
