package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/robbridges/webapp_v2/controllers"
	"github.com/robbridges/webapp_v2/migrations"
	"github.com/robbridges/webapp_v2/models"
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

	//user templates
	homeTpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	contactTpl := views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	faqTpl := views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	healthTpl := views.Must(views.ParseFS(templates.FS, "healthcheck.gohtml", "tailwind.gohtml"))
	signupTpl := views.Must(views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"))
	signInTpl := views.Must(views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"))
	currentUserTpl := views.Must(views.ParseFS(templates.FS, "currentuser.gohtml", "tailwind.gohtml"))
	forgotPasswordTpl := views.Must(views.ParseFS(templates.FS, "forgot_password.gohtml", "tailwind.gohtml"))
	checkEmailTpl := views.Must(views.ParseFS(templates.FS, "checkyouremail.gohtml", "tailwind.gohtml"))
	resetPwTpl := views.Must(views.ParseFS(templates.FS, "reset-password.gohtml", "tailwind.gohtml"))

	//gallery templates
	galleriesNewTpl := views.Must(views.ParseFS(templates.FS, "galleries/new.gohtml", "tailwind.gohtml"))
	gallertiesEditTpl := views.Must(views.ParseFS(templates.FS, "galleries/edit.gohtml", "tailwind.gohtml"))
	galleriesIndexTpl := views.Must(views.ParseFS(templates.FS, "galleries/index.gohtml", "tailwind.gohtml"))
	galleriesShowTpl := views.Must(views.ParseFS(templates.FS, "galleries/show.gohtml", "tailwind.gohtml"))

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
	pwResetService := models.PasswordResetService{
		DB: db,
	}

	smtpConfig := models.DefaultSMTPConfig()
	emailService := models.NewEmailService(smtpConfig)

	galleryService := &models.GalleryService{
		DB: db,
	}

	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
	}

	logger := &models.DBLogger{
		DB: db,
	}

	usersC := controllers.Users{
		UserService:          &userService,
		SessionService:       &sessionService,
		EmailService:         *emailService,
		PasswordResetService: pwResetService,
	}

	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	//user routes
	usersC.Templates.New = signupTpl
	usersC.Templates.SignIn = signInTpl
	usersC.Templates.CurrentUser = currentUserTpl
	usersC.Templates.ForgotPassword = forgotPasswordTpl
	usersC.Templates.CheckYourEmail = checkEmailTpl
	usersC.Templates.ResetPassword = resetPwTpl

	// gallery routes
	galleriesC.Templates.New = galleriesNewTpl
	galleriesC.Templates.Edit = gallertiesEditTpl
	galleriesC.Templates.Index = galleriesIndexTpl
	galleriesC.Templates.Show = galleriesShowTpl

	csrfKey := viper.GetString("CSRF_KEY")
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(viper.GetBool("CSRF_SECURE")), csrf.Path("/"))
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
	r.Get("/signin", usersC.SignIn)
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/signup", usersC.Create)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)
	r.Route("/currentuser", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.Route("/galleries", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/new", galleriesC.New)
			r.Get("/{id}/edit", galleriesC.Edit)
			r.Get("/", galleriesC.Index)
			r.Post("/", galleriesC.Create)
			r.Post("/{id}", galleriesC.Update)
			r.Post("/{id}/delete", galleriesC.Delete)
			r.Post("/{id}/images/{filename}/delete", galleriesC.DeleteImage)
		})
		r.Get("/{id}", galleriesC.Show)
		r.Get("/{id}/images/{filename}", galleriesC.Images)
	})
	r.NotFound(notFound)

	svr.ListenAndServe()
}
