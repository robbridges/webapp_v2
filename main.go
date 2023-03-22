package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/robbridges/webapp_v2/views"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func executeTemplate(w http.ResponseWriter, filepath string) {
	tpl, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("Parsing template %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
	viewTpl := views.Template{
		HtmlTpl: tpl,
	}
	viewTpl.Execute(w, nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "home.gohtml")
	executeTemplate(w, tplPath)
}

func contact(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "contact.gohtml")
	executeTemplate(w, tplPath)
}

func faq(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "faq.gohtml")
	executeTemplate(w, tplPath)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<p> Oh no, nothing's here..")
}

func main() {

	r := chi.NewRouter()
	svr := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	r.NotFound(notFound)
	r.Get("/", home)
	r.Get("/contact", contact)
	r.Get("/faq", faq)
	svr.ListenAndServe()
}
