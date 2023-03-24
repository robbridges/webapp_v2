package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/robbridges/webapp_v2/controllers"
	"github.com/robbridges/webapp_v2/templates"
	"github.com/robbridges/webapp_v2/views"
	"net/http"
)

func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "<p> Oh no, nothing's here..")
}

func main() {
	r := chi.NewRouter()
	homeTpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	contactTpl := views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	faqTpl := views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	healthTpl := views.Must(views.ParseFS(templates.FS, "healthcheck.gohtml", "tailwind.gohtml"))

	svr := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	r.NotFound(notFound)
	r.Get("/", controllers.StaticHandler(homeTpl))
	r.Get("/contact", controllers.StaticHandler(contactTpl))
	r.Get("/faq", controllers.FAQ(faqTpl))
	r.Get("/healthcheck", controllers.StaticHandler(healthTpl))
	svr.ListenAndServe()
}
