package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1> Welcome to LensLocked </h1>")
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Contact us at ...")
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
	svr.ListenAndServe()
}
