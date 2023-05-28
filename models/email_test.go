package models

import "testing"

func TestNewEmailService(t *testing.T) {
	mockConfig := SMTPConfig{
		HOST:     "testhost",
		Port:     1,
		Username: "rob",
		Password: "secret",
	}

	es := NewEmailService(mockConfig)
	if es.dialer == nil {
		t.Errorf("Dialer should not be nil")
	}
}
