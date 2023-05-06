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

	mockUserService := &models.MockUserService{}
	user := &models.User{ID: 1}
	mockUserService.AuthenticateFunc = func(email string, password string) (*models.User, error) {
		if password == "secure" {
			return user, nil
		} else {
			return nil, errors.New("invalid credentials")
		}
	}

	mockSessionService := &models.MockSessionService{}
	session := &models.Session{UserID: user.ID, Token: "abc123"}
	mockSessionService.CreateFunc = func(userID int) (*models.Session, error) {
		return session, nil
	}

	users := Users{
		Templates: struct {
			New         Template
			SignIn      Template
			CurrentUser Template
		}{},
		UserService:    mockUserService,
		SessionService: mockSessionService,
	}

	t.Run("happy path", func(t *testing.T) {
		mockLogger := &models.MockLogger{}
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/signin", nil)

		email := "test@test.com"
		password := "secure"

		// Wrap the handler with the logger middleware
		handler := models.LoggerMiddleware(mockLogger)(http.HandlerFunc(users.ProcessSignIn))

		// Add logger to the request context
		ctx := context.WithValue(r.Context(), "logger", mockLogger)
		r = r.WithContext(ctx)

		// Set form values
		data := url.Values{}
		data.Set("email", email)
		data.Set("password", password)
		r.PostForm = data

		// Call the wrapped handler
		handler.ServeHTTP(w, r)

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
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mockLogger := &models.MockLogger{}
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/signin", nil)

		email := "test@test.com"
		password := "wrongpassword"

		// Wrap the handler with the logger middleware
		handler := models.LoggerMiddleware(mockLogger)(http.HandlerFunc(users.ProcessSignIn))

		// Add logger to the request context
		ctx := context.WithValue(r.Context(), "logger", mockLogger)
		r = r.WithContext(ctx)

		// Set form values
		data := url.Values{}
		data.Set("email", email)
		data.Set("password", password)
		r.PostForm = data

		// Call the wrapped handler
		handler.ServeHTTP(w, r)

		// Assert response
		if w.Code != http.StatusBadRequest {
			t.Errorf("unexpected status code: got %v, want %v", w.Code, http.StatusUnauthorized)
		}

		// Assert no session cookie was set
		cookie := w.Header().Get("Set-Cookie")
		if cookie != "" {
			t.Errorf("unexpected cookie set: got %v, want %v", cookie, "")
		}
		// There are 4 errors in the test when ran together
		if len(mockLogger.ErrorLog) != 2 {
			t.Fatalf("unexpected number of errors logged: got %v, want %v", len(mockLogger.ErrorLog), 1)
		}
		err := mockLogger.ErrorLog[0].Error()
		expectedErr := "invalid credentials"
		if !strings.Contains(err, expectedErr) {
			t.Errorf("unexpected error message: got %v, want %v", err, expectedErr)
		}
	})

	t.Run("invalid session", func(t *testing.T) {
		mockLogger := &models.MockLogger{}

		mockUserService := &models.MockUserService{}
		mockSessionService := &models.MockSessionService{}

		users := Users{
			Templates: struct {
				New         Template
				SignIn      Template
				CurrentUser Template
			}{},
			UserService:    mockUserService,
			SessionService: mockSessionService,
		}

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/current_user", nil)

		// Wrap the handler with the logger middleware
		handler := models.LoggerMiddleware(mockLogger)(http.HandlerFunc(users.CurrentUser))

		// Add logger to the request context
		ctx := context.WithValue(r.Context(), "logger", mockLogger)
		r = r.WithContext(ctx)

		// Set session cookie
		cookie := new(http.Cookie)
		cookie.Name = CookieSession
		cookie.Value = "invalid session"
		cookie.Path = "/"
		http.SetCookie(w, cookie)

		// Call the wrapped handler
		handler.ServeHTTP(w, r)

		// Assert response
		if w.Code != http.StatusFound {
			t.Errorf("unexpected status code: got %v, want %v", w.Code, http.StatusUnauthorized)
		}

		// Assert no user was returned
		userID, ok := r.Context().Value("userID").(int)
		if ok {
			t.Errorf("unexpected userID in request context: got %v, want nil", userID)
		}

		// Assert error was logged
		if len(mockLogger.ErrorLog) != 2 {
			t.Fatalf("unexpected number of errors logged: got %v, want %v", len(mockLogger.ErrorLog), 1)
		}
		err := mockLogger.ErrorLog[0].Error()
		expectedErr := "http: named cookie not present"
		if !strings.Contains(err, expectedErr) {
			t.Errorf("unexpected error message: got %v, want %v", err, expectedErr)
		}

		// Assert cookie was deleted
		if len(w.Header()["Set-Cookie"]) == 0 {
			t.Errorf("cookie not deleted: got no cookie")
		}
	})

}
