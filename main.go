package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/robbridges/webapp_v2/controllers"
	"github.com/robbridges/webapp_v2/migrations"
	"github.com/robbridges/webapp_v2/models"
	"github.com/robbridges/webapp_v2/rand"
	"github.com/robbridges/webapp_v2/templates"
	"github.com/robbridges/webapp_v2/views"
	"github.com/spf13/viper"
	"net/http"
)

func init() {
	viper.SetConfigFile("local.env")
	viper.AddConfigPath("./")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("init: %w", err))
	}
}

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
	signInTpl := views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	currentUserTpl := views.Must(views.ParseFS(templates.FS, "currentuser.gohtml", "tailwind.gohtml"))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic("Error migrating, app closing")
	}

	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	logger := &models.DBLogger{
		DB: db,
	}

	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	usersC.Templates.New = signupTpl
	usersC.Templates.SignIn = signInTpl
	usersC.Templates.CurrentUser = currentUserTpl

	csrfKey := rand.GenerateRandByteSlice()
	csrfMw := csrf.Protect(csrfKey, csrf.Secure(false))
	svr := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	r.Use(models.LoggerMiddleware(logger))
	r.Use(csrfMw)
	r.Use(umw.SetUser)

	r.Get("/", controllers.StaticHandler(homeTpl))
	r.Get("/contact", controllers.StaticHandler(contactTpl))
	r.Get("/faq", controllers.FAQ(faqTpl))
	r.Get("/healthcheck", controllers.StaticHandler(healthTpl))
	r.Get("/signup", usersC.New)
	r.Post("/signup", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Route("/currentuser", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.NotFound(notFound)

	svr.ListenAndServe()
}
