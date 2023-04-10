package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/robbridges/webapp_v2/controllers"
	"github.com/robbridges/webapp_v2/models"
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
	signupTpl := views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	signinTpl := views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}

	usersC := controllers.Users{
		UserService: &userService,
	}
	usersC.Templates.New = signupTpl
	usersC.Templates.SignIn = signinTpl
	csrfKey := models.GenerateRandByteSlice()
	csrfMw := csrf.Protect(csrfKey, csrf.Secure(false))
	svr := http.Server{
		Addr:    ":8080",
		Handler: csrfMw(r),
	}

	r.Get("/", controllers.StaticHandler(homeTpl))
	r.Get("/contact", controllers.StaticHandler(contactTpl))
	r.Get("/faq", controllers.FAQ(faqTpl))
	r.Get("/healthcheck", controllers.StaticHandler(healthTpl))
	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Get("/currentuser", usersC.CurrentUser)
	r.NotFound(notFound)
	svr.ListenAndServe()
}
