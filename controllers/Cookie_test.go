package controllers

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
)

type mockLogger struct{}

func (mockLogger *mockLogger) Create(error) error {
	return nil
}

func TestSetCookie(t *testing.T) {
	// Create a new ResponseRecorder to capture the output.
	rr := httptest.NewRecorder()

	// Call the setCookie function with the mocked response and values.
	setCookie(rr, "session", "12345")

	// Check the cookies set in the response.
	cookies := rr.Result().Cookies()
	if len(cookies) != 1 {
		t.Errorf("expected one cookie, but got %d", len(cookies))
	}
	if cookies[0].Name != "session" {
		t.Errorf("expected cookie name to be session, but got %s", cookies[0].Name)
	}
	if cookies[0].Value != "12345" {
		t.Errorf("expected cookie value to be 12345, but got %s", cookies[0].Value)
	}

}

func TestReadCookie(t *testing.T) {
	rr := httptest.NewRecorder()

	setCookie(rr, "session", "12345")

	// Create a new Request using the captured response.
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Cookie", rr.Header().Get("Set-Cookie"))

	// Initialize a mock DBLogger and add it to the request context.
	logger := &mockLogger{}
	ctx := context.WithValue(req.Context(), "logger", logger)
	req = req.WithContext(ctx)

	// Call the readCookie function with the mocked request.
	value, err := readCookie(req, "session")

	// Check that the cookie value and error are correct.
	if err != nil {
		t.Errorf("unexpected error when reading cookie: %v", err)
	}
	if value != "12345" {
		t.Errorf("expected cookie value to be 12345, but got %s", value)
	}

}

func TestSetCookieAndDeleteCookie(t *testing.T) {
	// Create a mock response writer
	mockResponseWriter := httptest.NewRecorder()

	// Call the function to set a cookie
	cookieName := "myCookie"
	cookieValue := "myCookieValue"
	newCookie(cookieName, cookieValue)
	setCookie(mockResponseWriter, cookieName, cookieValue)

	// Check that the Set-Cookie header was set correctly
	resultingHeaders := mockResponseWriter.Header()
	setCookieHeaderValues, ok := resultingHeaders["Set-Cookie"]
	if !ok {
		t.Error("Expected Set-Cookie header to be set")
	} else if len(setCookieHeaderValues) != 1 {
		t.Errorf("Expected 1 Set-Cookie header value, got %d", len(setCookieHeaderValues))
	} else {
		expectedHeaderValue := fmt.Sprintf("%s=%s; Path=/; HttpOnly", cookieName, cookieValue)
		if setCookieHeaderValues[0] != expectedHeaderValue {
			t.Errorf("Expected Set-Cookie header value to be %q, got %q", expectedHeaderValue, setCookieHeaderValues[0])
		}
	}

	// Call the function to delete the cookie
	deleteCookie(mockResponseWriter, cookieName)

	// Check that the Set-Cookie header was set correctly
	resultingHeaders = mockResponseWriter.Header()
	setCookieHeaderValues, ok = resultingHeaders["Set-Cookie"]
	if !ok {
		t.Error("Expected Set-Cookie header to be set")
	} else if len(setCookieHeaderValues) != 2 {
		t.Errorf("Expected 1 Set-Cookie header value, got %d", len(setCookieHeaderValues))
	} else {
		expectedHeaderValue := fmt.Sprintf("%s=myCookieValue; Path=/; HttpOnly", cookieName)
		if setCookieHeaderValues[0] != expectedHeaderValue {
			t.Errorf("Expected Set-Cookie header value to be %q, got %q", expectedHeaderValue, setCookieHeaderValues[0])
		}
	}
}
