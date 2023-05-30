package integration_tests

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

	emailService := models.NewEmailService(cfg)

	testEmail := models.Email{
		To:        "admin@gallery.com",
		Subject:   "Test Email",
		Plaintext: "This is the plaintext content",
		HTML:      "<p>This is the HTML content</p>",
	}

	err := emailService.SendEmail(testEmail)

	if err != nil {
		t.Fatalf("Error sending email: %v", err)
	}
}
