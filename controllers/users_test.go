package controllers

import (
	"context"
	"errors"
	"github.com/robbridges/webapp_v2/models"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	// create mock template
	mockTemplate := &MockTemplate{
		ExecuteFunc: func(w http.ResponseWriter, r *http.Request, data interface{}) {
			// verify that the data contains the expected email
			if d, ok := data.(struct{ Email string }); ok {
				if d.Email != "test@example.com" {
					t.Errorf("Expected email to be 'test@example.com', got '%s'", d.Email)
				}
			} else {
				t.Error("Failed to cast data to expected struct")
			}
		},
	}

	// create instance of Users with mock dependencies
	users := Users{
		Templates: struct {
			New         Template
			SignIn      Template
			CurrentUser Template
		}{
			New: mockTemplate,
		},
	}

	// create mock http.ResponseWriter
	recorder := httptest.NewRecorder()

	// create mock http.Request with email form value
	req, _ := http.NewRequest("GET", "/?email=test@example.com", nil)

	// call the handler function
	users.New(recorder, req)
}

func TestSignIn(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/signin", nil)
	r.Form = url.Values{}
	r.Form.Add("email", "test@test.com")

	mockSignInTemplate := &MockTemplate{}
	mockSignInTemplate.ExecuteFunc = func(w http.ResponseWriter, r *http.Request, data interface{}) {}

	u := Users{
		Templates: struct {
			New         Template
			SignIn      Template
			CurrentUser Template
		}{
			SignIn: mockSignInTemplate,
		},
	}

	u.SignIn(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", w.Code, http.StatusOK)
	}
}

func TestUsers_ProcessSignIn(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/signin", nil)

	email := "test@test.com"
	password := "secure"

	// Create mock logger
	mockLogger := &models.MockLogger{}

	// Create mock user service
	mockUserService := &models.MockUserService{}
	user := &models.User{ID: 1}
	mockUserService.AuthenticateFunc = func(email string, password string) (*models.User, error) {
		if password == "secure" {
			return user, nil
		} else {
			return nil, errors.New("invalid credentials")
		}
	}

	// Create mock session service
	mockSessionService := &models.MockSessionService{}
	session := &models.Session{UserID: user.ID, Token: "abc123"}
	mockSessionService.CreateFunc = func(userID int) (*models.Session, error) {
		return session, nil
	}

	// Create Users struct with mocks
	users := Users{
		Templates: struct {
			New         Template
			SignIn      Template
			CurrentUser Template
		}{},
		UserService:    mockUserService,
		SessionService: mockSessionService,
	}

	// Add logger middleware to request context
	ctx := context.WithValue(r.Context(), "logger", mockLogger)
	r = r.WithContext(ctx)

	// Set form values
	data := url.Values{}
	data.Set("email", email)
	data.Set("password", password)
	r.PostForm = data

	// Call ProcessSignIn function
	users.ProcessSignIn(w, r)

	// Assert response
	if w.Code != http.StatusFound {
		t.Errorf("unexpected status code: got %v, want %v", w.Code, http.StatusFound)
	}

	// Assert session cookie was set
	cookie := w.Header().Get("Set-Cookie")
	if !strings.Contains(cookie, CookieSession+"="+session.Token) {
		t.Errorf("cookie not set correctly: got %v, want %v", cookie, CookieSession+"="+session.Token)
	}

	// Assert no error was logged
	if len(mockLogger.ErrorLog) > 0 {
		t.Errorf("unexpected error logged: %v", mockLogger.ErrorLog[0])
	}
}
