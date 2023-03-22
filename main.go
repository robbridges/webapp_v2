package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1> Welcome to LensLocked! </h1>")
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>ContactPage</h1><p>To get in touch, email me at"+
		"<a href=\"mailto:admin@lenslocked.com\">admin@lenslocked.com</a>")
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1> FAQ PAGE </h1>"+
		"<p> Q: Is there a free version?</p> \n "+
		"<p> A: Yes! We offfer a free 30 day trial of all paid plans </p>")
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
