package models

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
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
	err = os.RemoveAll(svc.galleryDir(gallery.ID))
	if err != nil {
		return fmt.Errorf("delete gallery images: %w", err)
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
				GalleryID: galleryID,
				Path:      file,
				Filename:  filepath.Base(file),
			})
		}
	}
	return images, nil
}

func (svc *GalleryService) CreateImage(galleryID int, filename string, contents io.ReadSeeker) error {
	err := CheckContentType(contents, svc.imageContentTypes())
	if err != nil {
		return fmt.Errorf("creating image: %v: %v", filename, err)
	}

	err = checkExtension(filename, svc.extensions())
	if err != nil {
		return fmt.Errorf("creating image: %v: %v", filename, err)
	}

	galleryDir := svc.galleryDir(galleryID)
	err = os.MkdirAll(galleryDir, 0755)
	if err != nil {
		return fmt.Errorf("create gallery-%d directory: %v", galleryID, err)
	}
	imagePath := filepath.Join(galleryDir, filename)
	dst, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("creating image file: %v", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, contents)
	if err != nil {
		return fmt.Errorf("copying contents to image: %v", err)
	}

	return nil
}

func (svc *GalleryService) Image(galleryID int, filename string) (Image, error) {
	imagePath := filepath.Join(svc.galleryDir(galleryID), filename)
	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Image{}, ErrNoData
		}
		return Image{}, fmt.Errorf("querying for image: %v", err)
	}

	return Image{
		Filename:  filename,
		GalleryID: galleryID,
		Path:      imagePath,
	}, nil
}

func (svc *GalleryService) DeleteImage(galleryID int, filename string) error {
	image, err := svc.Image(galleryID, filename)
	if err != nil {
		return fmt.Errorf("deleting image: %v", err)
	}

	err = os.Remove(image.Path)
	if err != nil {
		return fmt.Errorf("deleting image %v", err)
	}
	return nil
}

func (svc *GalleryService) galleryDir(id int) string {
	imagesDir := svc.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}

func (svc *GalleryService) extensions() []string {
	return []string{".png", ".jpg", ".jpeg", ".gif"}
}

func (svc *GalleryService) imageContentTypes() []string {
	return []string{"image/png", "image/jpeg", "image/gif"}
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
