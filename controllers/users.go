package controllers

import (
	"fmt"
	"github.com/robbridges/webapp_v2/models"
	"net/http"
)

type Users struct {
	Templates struct {
		New         Template
		SignIn      Template
		CurrentUser Template
	}
	UserService    *models.UserService
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
	logger := r.Context().Value("logger").(*models.DBLogger)
	var data struct {
		Email    string
		Password string
	}

	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		logger.Create(err)
		http.Error(w, "Something went wrong", http.StatusBadRequest)
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		logger.Create(err)
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	setCookie(w, CookieSession, session.Token)

	http.Redirect(w, r, "/currentuser", http.StatusFound)

}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*models.DBLogger)
	tokenCookie, err := readCookie(r, CookieSession)
	if err != nil {
		logger.Create(err)
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	user, err := u.SessionService.User(tokenCookie)
	if err != nil {
		logger.Create(err)
		fmt.Println(err)
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
	logger := r.Context().Value("logger").(*models.DBLogger)
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		logger.Create(err)
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		logger.Create(err)
		fmt.Println(err)
		//TODO: Long term there's a better way to handle this without a confusing redirect
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/currentuser", http.StatusFound)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(*models.DBLogger)
	token, err := readCookie(r, CookieSession)
	if err != nil {
		logger.Create(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	err = u.SessionService.DeleteSession(token)
	if err != nil {
		logger.Create(err)
		fmt.Errorf("delete session %w", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}
