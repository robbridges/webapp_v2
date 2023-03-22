package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	tplPath := filepath.Join("templates", "home.gohtml")
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		log.Printf("Parsing template %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
	err = tpl.Execute(w, nil)
	if err != nil {
		log.Printf("Executing template %v", err)
		http.Error(w, "There was an error executing the template.", http.StatusInternalServerError)
		return
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>ContactPage</h1><p>To get in touch, email me at "+
		"<a href=\"mailto:admin@lenslocked.com\">admin@lenslocked.com</a>")
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1> FAQ PAGE </h1>"+
		"<p> Q: Is there a free version?</p> \n "+
		"<p> A: Yes! We offer a free 30 day trial of all paid plans </p>",
	)
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
