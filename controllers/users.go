package controllers

import (
	"fmt"
	"github.com/robbridges/webapp_v2/models"
	"net/http"
)

type Users struct {
	Templates struct {
		New Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")

	u.Templates.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "email:", r.FormValue("email"))
	fmt.Fprintf(w, "password:", r.FormValue("password"))
	fmt.Fprintf(w, "file:", r.FormValue("file"))
}
