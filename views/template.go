package views

import (
	"html/template"
	"log"
	"net/http"
)

type Template struct {
	HtmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")

	err := t.HtmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("Executing template %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}
