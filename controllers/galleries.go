package controllers

import (
	"fmt"
	"github.com/robbridges/webapp_v2/context"
	"github.com/robbridges/webapp_v2/models"
	"net/http"
)

type Galleries struct {
	Templates struct {
		New Template
	}
	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Templates.New.Execute(w, r, data)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(models.LogInterface)

	var data struct {
		UserID int
		Title  string
	}

	data.UserID = context.User(r.Context()).ID
	data.Title = r.FormValue("title")
	gallery, err := g.GalleryService.Create(data.Title, data.UserID)
	if err != nil {
		logger.Create(err)
		g.Templates.New.Execute(w, r, data, err)
		return
	}

	//TODO implement this page
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}
