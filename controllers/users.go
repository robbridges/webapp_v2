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
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "User Created: %+v", user)
}
