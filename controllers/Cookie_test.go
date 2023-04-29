package controllers

import (
	"net/http/httptest"
	"testing"
)

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
