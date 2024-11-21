package providers

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/loureirovinicius/cleanup/helpers/logger"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/mock"
)

// Mock AWS client and services
type mockAWSConfigLoader struct {
	mock.Mock
}

func (m *mockAWSConfigLoader) LoadDefaultConfig(ctx context.Context, opts ...func(*config.LoadOptions) error) (aws.Config, error) {
	args := m.Called(ctx)
	return args.Get(0).(aws.Config), args.Error(1)
}

func TestLoadProvider(t *testing.T) {
	// Mock dependencies
	ctx := context.Background()

	mockLoader := &mockAWSConfigLoader{}

	// Mock AWS config
	awsConfig := aws.Config{}
	mockLoader.On("LoadDefaultConfig", ctx).Return(awsConfig, nil)

	// Initialize logger
	logger.InitializeLogger("info", "text", os.Stdout)

	// Test case: Unsupported provider
	t.Run("Unsupported Provider", func(t *testing.T) {
		_, err := LoadProvider(ctx, "gcp", "ebs")
		if err == nil || err.Error() != "provider gcp is not supported" {
			t.Errorf("Expected unsupported provider error, got: %v", err)
		}
	})

	// Test case: AWS provider with supported service
	t.Run("AWS Provider with Supported Service", func(t *testing.T) {
		viper.Set("aws.region", "us-east-1") // Mock configuration
		provider := AWS{}

		// Create a valid aws.Config
		client := &awsConfig

		service, err := provider.loadService(ctx, client, "ebs")
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if service == nil {
			t.Error("Expected valid service, got nil")
		}
	})
}

func TestAWSInitialize(t *testing.T) {
	provider := AWS{}

	// Test case: Missing region
	t.Run("Missing Region", func(t *testing.T) {
		viper.Reset() // Clear mocked config
		err := provider.loadConfig()
		if err == nil || err.Error() != "AWS region can't be empty" {
			t.Errorf("Expected missing region error, got: %v", err)
		}
	})

	// Test case: Valid configuration
	t.Run("Valid Configuration", func(t *testing.T) {
		viper.Set("aws.region", "us-east-1")
		err := provider.loadConfig()
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})
}
