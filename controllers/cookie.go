package controllers

import (
	"fmt"
	"github.com/robbridges/webapp_v2/models"
	"net/http"
)

const (
	CookieSession = "session"
)

func newCookie(name, value string) *http.Cookie {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
	}

	return &cookie
}

func setCookie(w http.ResponseWriter, name, value string) {
	cookie := newCookie(name, value)
	http.SetCookie(w, cookie)
}

func readCookie(r *http.Request, name string) (string, error) {
	logger := r.Context().Value("logger").(models.LogInterface)
	c, err := r.Cookie(name)
	if err != nil {
		logger.Create(err)
		return "", fmt.Errorf("cookie %s: read error: %w", name, err)
	}
	return c.Value, nil
}

func deleteCookie(w http.ResponseWriter, name string) {
	cookie := newCookie(name, "")
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)
}
