package main

import (
	"fmt"
	"html/template"
	"net/http"
)

const portNumber = ":8080"

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

func main() {

	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	fmt.Println(fmt.Sprintf("Starting application on port: %s", portNumber))
	_ = http.ListenAndServe(portNumber, nil)
}
