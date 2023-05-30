package library

import (
	"github.com/robbridges/webapp_v2/models"
	"github.com/spf13/viper"
	"testing"
)

func TestSendEmail(t *testing.T) {
	viper.SetConfigFile("../local.env")
	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	cfg := models.SMTPConfig{
		HOST:     viper.GetString("EMAIL_HOST"),
		Port:     viper.GetInt("EMAIL_PORT"),
		Username: viper.GetString("EMAIL_USERNAME"),
		Password: viper.GetString("EMAIL_PASSWORD"),
	}

	// Create an instance of EmailService
	emailService := models.NewEmailService(cfg)

	// Create a test Email object
	testEmail := models.Email{
		To:        "admin@gallery.com",
		Subject:   "Test Email",
		Plaintext: "This is the plaintext content",
		HTML:      "<p>This is the HTML content</p>",
	}

	// Call the SendEmail function
	err := emailService.SendEmail(testEmail)

	// Perform assertions to verify the result
	if err != nil {
		t.Fatalf("Error sending email: %v", err)
	}
}
