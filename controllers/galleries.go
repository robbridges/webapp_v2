package controllers

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/robbridges/webapp_v2/context"
	"github.com/robbridges/webapp_v2/models"
	"math/rand"
	"net/http"
	"strconv"
)

type Galleries struct {
	Templates struct {
		New   Template
		Show  Template
		Edit  Template
		Index Template
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

	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	var data struct {
		ID     int
		Title  string
		Images []string
	}

	data.ID = gallery.ID
	data.Title = gallery.Title

	// fake images
	for i := 0; i < 20; i++ {
		w, h := rand.Intn(500)+200, rand.Intn(500)+200
		randomImage := fmt.Sprintf("https://placeimg.com/%d/%d/any", w, h)
		data.Images = append(data.Images, randomImage)
	}
	g.Templates.Show.Execute(w, r, data)

}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}

	var data struct {
		ID    int
		Title string
	}

	data.ID = gallery.ID
	data.Title = gallery.Title
	g.Templates.Edit.Execute(w, r, data)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(models.LogInterface)
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = userMustOwnGallery(w, r, gallery)
	if err != nil {
		return
	}

	gallery.Title = r.FormValue("title")
	err = g.GalleryService.Update(gallery)
	if err != nil {
		logger.Create(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(models.LogInterface)
	type Gallery struct {
		ID    int
		Title string
	}

	var data struct {
		Galleries []Gallery
	}
	user := context.User(r.Context())

	galleries, err := g.GalleryService.ByUserID(user.ID)
	if err != nil {
		logger.Create(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:    gallery.ID,
			Title: gallery.Title,
		})
	}

	g.Templates.Index.Execute(w, r, data)

}

func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value("logger").(models.LogInterface)
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = g.GalleryService.Delete(gallery)
	if err != nil {
		logger.Create(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

type galleryOpt func(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error

func (g Galleries) galleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOpt) (*models.Gallery, error) {
	logger := r.Context().Value("logger").(models.LogInterface)
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return nil, err
	}
	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNoData) {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return nil, err
		}
		logger.Create(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return nil, err
	}

	for _, opt := range opts {
		err = opt(w, r, gallery)
		if err != nil {
			return nil, err
		}
	}

	return gallery, nil
}

func userMustOwnGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized to edit this gallery", http.StatusForbidden)
		return fmt.Errorf("user does not have access to this gallery")
	}

	return nil
}
