package models

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB        *sql.DB
	ImagesDir string
}

func (svc *GalleryService) Create(title string, userID int) (*Gallery, error) {
	gallery := Gallery{
		Title:  title,
		UserID: userID,
	}
	row := svc.DB.QueryRow(`
		INSERT INTO galleries (title, user_id)
		VALUES ($1, $2) RETURNING id;`, gallery.Title, gallery.UserID)
	err := row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %v", err)
	}
	return &gallery, nil
}

func (svc *GalleryService) ByID(id int) (*Gallery, error) {
	gallery := Gallery{
		ID: id,
	}
	row := svc.DB.QueryRow(`
		SELECT title, user_id
		FROM galleries
		WHERE id = $1;`, gallery.ID)
	err := row.Scan(&gallery.Title, &gallery.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoData
		}
		return nil, fmt.Errorf("query gallery by id: %w", err)
	}
	return &gallery, nil
}

func (svc *GalleryService) ByUserID(userID int) ([]Gallery, error) {
	rows, err := svc.DB.Query(`
		SELECT id, title
		FROM galleries
		WHERE user_id = $1;`, userID)
	if err != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	var galleries []Gallery
	for rows.Next() {
		gallery := Gallery{
			UserID: userID,
		}
		err := rows.Scan(&gallery.ID, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("query galleries by user: %w", err)
		}
		galleries = append(galleries, gallery)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	return galleries, nil
}

func (svc *GalleryService) Update(gallery *Gallery) error {
	_, err := svc.DB.Exec(`
		UPDATE galleries
		SET title = $2
	 	WHERE ID =$1;`, gallery.ID, gallery.Title,
	)

	if err != nil {
		return fmt.Errorf("update gallery: %v", err)
	}

	return nil
}

func (svc *GalleryService) Delete(gallery *Gallery) error {
	_, err := svc.DB.Exec(`
	DELETE from galleries
	WHERE id = $1;`, gallery.ID,
	)
	if err != nil {
		return fmt.Errorf("delete gallery: %v", err)
	}
	return nil
}

func (svc *GalleryService) Images(galleryID int) ([]Image, error) {
	globPattern := filepath.Join(svc.galleryDir(galleryID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery images: %w", err)
	}
	var images []Image
	for _, file := range allFiles {
		if hasExtension(file, svc.extensions()) {
			images = append(images, Image{
				Path: file,
			})
		}
	}
	return images, nil
}

func (svc GalleryService) galleryDir(id int) string {
	imagesDir := svc.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}

func (svc *GalleryService) extensions() []string {
	return []string{".png", ".jpg", ".jpeg", ".gif"}
}

func hasExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}
