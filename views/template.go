package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Template struct {
	HtmlTpl *template.Template
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func Parse(filepath string) (Template, error) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing templage: %w", err)
	}
	return Template{
		HtmlTpl: tpl,
	}, nil
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
