package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/robbridges/webapp_v2/controllers"
	"github.com/robbridges/webapp_v2/views"
	"net/http"
	"path/filepath"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<p> Oh no, nothing's here..")
}

func main() {
	r := chi.NewRouter()
	hometpl, err := views.Parse(filepath.Join("templates", "home.gohtml"))
	if err != nil {
		panic(err)
	}
	contactTpl, err := views.Parse(filepath.Join("templates", "contact.gohtml"))
	if err != nil {
		panic(err)
	}
	faqTpl, err := views.Parse(filepath.Join("templates", "faq.gohtml"))
	if err != nil {
		panic(err)
	}

	svr := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	r.NotFound(notFound)
	r.Get("/", controllers.StaticHandler(hometpl))
	r.Get("/contact", controllers.StaticHandler(contactTpl))
	r.Get("/faq", controllers.StaticHandler(faqTpl))
	svr.ListenAndServe()
}
