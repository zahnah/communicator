package render

import (
	"fmt"
	"html/template"
	"net/http"
)

func Template(writer http.ResponseWriter, tmpl string) {
	parseTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	err := parseTemplate.Execute(writer, nil)
	if err != nil {
		fmt.Println("error parsing template:", err)
		return
	}
}
