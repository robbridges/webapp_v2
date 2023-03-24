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
	homeTpl := views.Must(views.Parse(filepath.Join("templates", "home.gohtml")))
	contactTpl := views.Must(views.Parse(filepath.Join("templates", "contact.gohtml")))
	faqTpl := views.Must(views.Parse(filepath.Join("templates", "faq.gohtml")))

	svr := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	r.NotFound(notFound)
	r.Get("/", controllers.StaticHandler(homeTpl))
	r.Get("/contact", controllers.StaticHandler(contactTpl))
	r.Get("/faq", controllers.StaticHandler(faqTpl))
	svr.ListenAndServe()
}
