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

	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/?email=test@example.com", nil)

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

		if w.Code != http.StatusFound {
			t.Errorf("unexpected status code: got %v, want %v", w.Code, http.StatusFound)
		}

		cookie := w.Header().Get("Set-Cookie")
		if !strings.Contains(cookie, CookieSession+"="+session.Token) {
			t.Errorf("cookie not set correctly: got %v, want %v", cookie, CookieSession+"="+session.Token)
		}

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

		handler := models.LoggerMiddleware(mockLogger)(http.HandlerFunc(users.ProcessSignIn))

		ctx := context.WithValue(r.Context(), "logger", mockLogger)
		r = r.WithContext(ctx)

		data := url.Values{}
		data.Set("email", email)
		data.Set("password", password)
		r.PostForm = data

		handler.ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("unexpected status code: got %v, want %v", w.Code, http.StatusUnauthorized)
		}

		cookie := w.Header().Get("Set-Cookie")
		if cookie != "" {
			t.Errorf("unexpected cookie set: got %v, want %v", cookie, "")
		}

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

		if len(w.Header()["Set-Cookie"]) == 0 {
			t.Errorf("cookie not deleted: got no cookie")
		}
	})
}

func TestUsers_CurrentUser(t *testing.T) {
	mockLogger := &models.MockLogger{}

	mockSessionService := &models.MockSessionService{}
	mockUserService := &models.MockUserService{}

	user := &models.User{
		ID:    1,
		Email: "test@test.com",
	}

	mockUserService.AuthenticateFunc = func(email string, password string) (*models.User, error) {
		if password == "secure" {
			return user, nil
		} else {
			return nil, errors.New("invalid credentials")
		}
	}
	mockSessionService.UserFunc = func(token string) (*models.User, error) {
		if token == "valid_token" {
			return user, nil
		} else {
			return nil, errors.New("session read error")
		}
	}

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
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/signin", nil)

		cookie := &http.Cookie{
			Name:  CookieSession,
			Value: "test_cookie_value",
		}
		r.AddCookie(cookie)

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

		setCookie(w, CookieSession, session.Token)

		// Call readCookie to read the "session" cookie value
		value, err := readCookie(r, CookieSession)

		// Check if the value and error are as expected
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if value != "test_cookie_value" {
			t.Errorf("unexpected cookie value: got %v, want %v", value, "test_cookie_value")
		}

		if len(mockLogger.ErrorLog) > 0 {
			t.Errorf("unexpected error logged: %v", mockLogger.ErrorLog[0])
		}
	})

	t.Run("invalid session", func(t *testing.T) {
		mockLogger := &models.MockLogger{}
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/current_user", nil)

		// Wrap the handler with the logger middleware
		handler := models.LoggerMiddleware(mockLogger)(http.HandlerFunc(users.CurrentUser))

		// Add logger to the request context
		ctx := context.WithValue(r.Context(), "logger", mockLogger)
		r = r.WithContext(ctx)

		// Set session cookie
		setCookie(w, CookieSession, "invalid session")

		// Call the wrapped handler
		handler.ServeHTTP(w, r)

		// Assert response
		if w.Code != http.StatusFound {
			t.Errorf("unexpected status code: got %v, want %v", w.Code, http.StatusFound)
		}

		// Assert no user was returned
		userID, ok := r.Context().Value("userID").(int)
		if ok {
			t.Errorf("unexpected userID in request context: got %v, want nil", userID)
		}

		// Assert error was logged
		if len(mockLogger.ErrorLog) != 2 {
			t.Fatalf("unexpected number of errors logged: got %v, want 1", len(mockLogger.ErrorLog))
		}
		err := mockLogger.ErrorLog[0].Error()
		expectedErr := "http: named cookie not present"
		if !strings.Contains(err, expectedErr) {
			t.Errorf("unexpected error message: got %v, want %v", err, expectedErr)
		}

		// Assert cookie was deleted
		cookies := r.Cookies()
		if len(cookies) != 0 {
			t.Errorf("unexpected number of cookies: got %v, want 0", len(cookies))
		}
	})
}

func TestUsers_Create(t *testing.T) {
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

	t.Run("happy path", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/create", nil)

		// Wrap the handler with the logger middleware
		handler := models.LoggerMiddleware(mockLogger)(http.HandlerFunc(users.Create))

		// Add logger to the request context
		ctx := context.WithValue(r.Context(), "logger", mockLogger)
		r = r.WithContext(ctx)

		// Set form values
		data := url.Values{}
		data.Set("email", "test@test.com")
		data.Set("password", "secure")
		r.PostForm = data

		// Mock user service Create method
		user := &models.User{
			ID:    1,
			Email: "test@test.com",
		}
		mockUserService.CreateFunc = func(email string, password string) (*models.User, error) {
			return user, nil
		}

		// Mock session service Create method
		session := &models.Session{
			UserID: user.ID,
			Token:  "test_session_token",
		}
		mockSessionService.CreateFunc = func(userID int) (*models.Session, error) {
			return session, nil
		}

		// Call the wrapped handler
		handler.ServeHTTP(w, r)

		// Assert response
		if w.Code != http.StatusFound {
			t.Errorf("unexpected status code: got %v, want %v", w.Code, http.StatusFound)
		}

		// Assert session cookie was set
		cookies := w.Result().Cookies()
		if len(cookies) != 1 {
			t.Errorf("unexpected number of cookies: got %v, want %v", len(cookies), 1)
		}
		if cookies[0].Name != CookieSession {
			t.Errorf("unexpected cookie name: got %v, want %v", cookies[0].Name, CookieSession)
		}
		if cookies[0].Value != session.Token {
			t.Errorf("unexpected cookie value: got %v, want %v", cookies[0].Value, session.Token)
		}
		if cookies[0].Path != "/" {
			t.Errorf("unexpected cookie path: got %v, want %v", cookies[0].Path, "/")
		}

		// Assert redirection
		if w.Header().Get("Location") != "/currentuser" {
			t.Errorf("unexpected redirect location: got %v, want %v", w.Header().Get("Location"), "/currentuser")
		}

		if len(mockLogger.ErrorLog) > 0 {
			t.Errorf("unexpected error logged: %v", mockLogger.ErrorLog[0])
		}
	})

	t.Run("create user error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/signup", nil)

		// Wrap the handler with the logger middleware
		handler := models.LoggerMiddleware(mockLogger)(http.HandlerFunc(users.Create))

		// Add logger to the request context
		ctx := context.WithValue(r.Context(), "logger", mockLogger)
		r = r.WithContext(ctx)

		// Set form values
		data := url.Values{}
		data.Set("email", "test@test.com")
		data.Set("password", "weakpassword")
		r.PostForm = data

		// Mock UserService to return an error when creating a user
		mockUserService.CreateFunc = func(email string, password string) (*models.User, error) {
			return nil, errors.New("create user error")
		}
		defer func() {
			mockUserService.CreateFunc = nil
		}()

		// Call the wrapped handler
		handler.ServeHTTP(w, r)

		// Assert response
		if w.Code != http.StatusInternalServerError {
			t.Errorf("unexpected status code: got %v, want %v", w.Code, http.StatusInternalServerError)
		}

		// Assert error was logged
		if len(mockLogger.ErrorLog) != 1 {
			t.Fatalf("unexpected number of errors logged: got %v, want 1", len(mockLogger.ErrorLog))
		}
		err := mockLogger.ErrorLog[0].Error()
		expectedErr := "create user error"
		if !strings.Contains(err, expectedErr) {
			t.Errorf("unexpected error message: got %v, want %v", err, expectedErr)
		}

	})
}

func TestUsers_ProcessSignOut(t *testing.T) {
	mockLogger := &models.MockLogger{}
	mockSessionService := &models.MockSessionService{}

	mockSessionService.DeleteSessionFunc = func(token string) error {
		if token == "valid_token" {
			return nil
		} else {
			return errors.New("session delete error")
		}
	}

	users := Users{
		Templates: struct {
			New         Template
			SignIn      Template
			CurrentUser Template
		}{},
		SessionService: mockSessionService,
	}

	handler := models.LoggerMiddleware(mockLogger)(http.HandlerFunc(users.ProcessSignOut))

	t.Run("happy path", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/signout", nil)

		cookie := &http.Cookie{
			Name:  CookieSession,
			Value: "valid_token",
		}
		r.AddCookie(cookie)

		ctx := context.WithValue(r.Context(), "logger", mockLogger)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)

		// Assert response
		if w.Code != http.StatusFound {
			t.Errorf("unexpected status code: got %v, want %v", w.Code, http.StatusFound)
		}

		deleteCookie(w, CookieSession)

		// Assert cookie was deleted
		cookies := r.Cookies()
		if len(cookies) != 1 {
			t.Errorf("unexpected number of cookies: got %d, want 1", len(cookies))
		} else {
			if cookies[0].Name != CookieSession {
				t.Errorf("unexpected cookie name: got %s, want %s", cookies[0].Name, CookieSession)
			}
			if cookies[0].MaxAge != 0 {
				t.Errorf("unexpected cookie max age: got %d, want -1", cookies[0].MaxAge)
			}
		}

		if len(mockLogger.ErrorLog) > 0 {
			t.Errorf("unexpected error logged: %v", mockLogger.ErrorLog[0])
		}
	})

	t.Run("sad path", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/signout", nil)

		cookie := &http.Cookie{
			Name:  CookieSession,
			Value: "invalid_token",
		}
		r.AddCookie(cookie)

		ctx := context.WithValue(r.Context(), "logger", mockLogger)
		r = r.WithContext(ctx)

		handler.ServeHTTP(w, r)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("unexpected status code: got %v, want %v", w.Code, http.StatusInternalServerError)
		}

		// Check if the cookie was deleted
		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Name == CookieSession {
				t.Errorf("cookie not deleted: %v", cookie)
			}
		}

		if len(mockLogger.ErrorLog) == 0 {
			t.Error("expected error not logged")
		}
	})
}
