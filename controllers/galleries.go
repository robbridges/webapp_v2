package controllers

import (
	"github.com/robbridges/webapp_v2/models"
	"net/http"
)

type Galleries struct {
	Templates struct {
		new Template
	}
	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Templates.new.Execute(w, r, data)
}
