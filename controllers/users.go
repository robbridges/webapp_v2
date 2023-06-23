package controllers

import (
	"errors"
	"fmt"
	"github.com/robbridges/webapp_v2/context"
	puberror "github.com/robbridges/webapp_v2/errors"
	"github.com/robbridges/webapp_v2/models"
	"net/http"
	"net/url"
)

type Users struct {
	Templates struct {
		New            Template
		SignIn         Template
		CurrentUser    Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassword  Template
	}
	UserService          models.UserServiceInterface
	SessionService       models.SessionServiceInterface
	PasswordResetService models.PasswordResetService
	EmailService         models.EmailService
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")

	u.Templates.New.Execute(w, r, data)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(models.LogInterface)
	var data struct {
		Email    string
		Password string
	}

	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		err = logger.Create(err)
		if err != nil {
			fmt.Println("Err creating log ")
		}
		http.Error(w, "Something went wrong", http.StatusBadRequest)
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		err = logger.Create(err)
		if err != nil {
			fmt.Println("Err creating log ")
		}
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, session.Token)

	http.Redirect(w, r, "/galleries", http.StatusFound)

}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	var data struct {
		Email string
	}

	data.Email = user.Email

	u.Templates.CurrentUser.Execute(w, r, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}

	logger := r.Context().Value("logger").(models.LogInterface)
	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Create(data.Email, data.Password)
	if err != nil {
		logger.Create(err)

		if errors.Is(err, models.ErrEmailTaken) {
			err = puberror.Public(err, "That email is already taken")
		}
		u.Templates.New.Execute(w, r, data, err)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		err = logger.Create(err)
		if err != nil {
			fmt.Println("Err creating log ")
		}
		fmt.Println(err)
		//TODO: Long term there's a better way to handle this without a confusing redirect
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(models.LogInterface)
	token, err := readCookie(r, CookieSession)
	if err != nil {
		err = logger.Create(err)
		if err != nil {
			fmt.Println("Err creating log ")
		}
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	err = u.SessionService.DeleteSession(token)
	if err != nil {
		err = logger.Create(err)
		if err != nil {
			fmt.Println("Err creating log ")
		}

		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := readCookie(r, CookieSession)
		if err != nil {

			next.ServeHTTP(w, r)
			return
		}

		user, err := umw.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := r.Context()

		ctx = context.WithUser(ctx, user)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(models.LogInterface)
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		if err.Error() == models.ErrNoData.Error() {
			err = puberror.Public(err, "That email is already taken")
			u.Templates.ForgotPassword.Execute(w, r, data, err)
			return
		}

		fmt.Println(err)
		http.Error(w, "Something went wrong, check your credentials", http.StatusInternalServerError)
		return
	}

	vals := url.Values{
		"token": {pwReset.Token},
	}

	// we actually should let this block the account instead of doing it in a go routine, a user can't move forward until
	// email has been sent, thought was put into it to make this a backround job in a go routine, but it risks
	// the user getting a "Check your email" http response before the email was sent if that is behind for whatever
	// reason. A welcome email would be different.

	resetUrl := "https://www.webgallery.com/reset-pw?" + vals.Encode()
	err = u.EmailService.ForgotPassword(data.Email, resetUrl)
	if err != nil {
		err = logger.Create(err)
		if err != nil {
			fmt.Println("Err creating log ")
		}
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	// don't render token her, user needs to confirm email account access
	u.Templates.CheckYourEmail.Execute(w, r, data)
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}

	data.Token = r.FormValue("token")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(models.LogInterface)
	var data struct {
		Token    string
		Password string
	}

	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")

	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		logger.Create(err)
		http.Error(w, "Somewent went wrong", http.StatusInternalServerError)
		return
	}

	err = u.UserService.UpdatePassword(user.ID, data.Password)

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		logger.Create(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
	}

	setCookie(w, CookieSession, session.Token)

	http.Redirect(w, r, "/currentuser", http.StatusFound)
}
