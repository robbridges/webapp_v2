package models

import (
	"database/sql"
	"fmt"
)

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB *sql.DB
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
		return nil, fmt.Errorf("query gallery by id: %w", err)
	}
	return &gallery, nil
}
