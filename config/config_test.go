// config_test.go
package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestStart(t *testing.T) {
	// Test case 1: Valid provider with missing config.yaml
	t.Run("Valid AWS provider, missing config.yaml", func(t *testing.T) {
		// Cleanup environment
		defer viper.Reset()

		err := Start("aws")
		if err == nil || err.Error() != "'config.yaml' file does not exist" {
			t.Errorf("Expected missing file error, got: %v", err)
		}
	})

	// Test case 2: Unsupported provider
	t.Run("Unsupported provider", func(t *testing.T) {
		err := Start("gcp")
		if err == nil || err.Error() != "provider gcp is not supported" {
			t.Errorf("Expected unsupported provider error, got: %v", err)
		}
	})

	// Test case 3: Valid provider with valid config.yaml
	t.Run("Valid AWS provider, valid config.yaml", func(t *testing.T) {
		// Cleanup environment
		defer viper.Reset()

		// Create a temporary config.yaml file
		file, err := os.Create("config.yaml")
		if err != nil {
			t.Fatalf("Failed to create temporary config.yaml: %v", err)
		}
		defer os.Remove(file.Name())
		_, _ = file.WriteString(`region: "us-east-1"`)

		err = Start("aws")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})
}

func TestLoadAWSConfig(t *testing.T) {
	// Cleanup environment
	defer viper.Reset()

	// Mock environment variables
	_ = os.Setenv("AWS_REGION", "us-east-1")
	_ = os.Setenv("AWS_PROFILE_NAME", "default")
	_ = os.Setenv("AWS_ACCESS_KEY", "mock-access-key")
	_ = os.Setenv("AWS_SECRET_KEY", "mock-secret-key")

	err := loadAWSConfig()
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Check if variables are properly set
	if viper.GetString("region") != "us-east-1" {
		t.Errorf("Expected region to be 'us-east-1', got: %s", viper.GetString("region"))
	}
	if viper.GetString("profile.name") != "default" {
		t.Errorf("Expected profile.name to be 'default', got: %s", viper.GetString("profile.name"))
	}
	if viper.GetString("credentials.access_key") != "mock-access-key" {
		t.Errorf("Expected credentials.access_key to be 'mock-access-key', got: %s", viper.GetString("credentials.access_key"))
	}
	if viper.GetString("credentials.secret_key") != "mock-secret-key" {
		t.Errorf("Expected credentials.secret_key to be 'mock-secret-key', got: %s", viper.GetString("credentials.secret_key"))
	}
}
