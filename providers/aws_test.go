package providers

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/loureirovinicius/cleanup/helpers/logger"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadAWSService(t *testing.T) {
	// Mock dependencies
	ctx := context.Background()

	provider := AWS{}

	// Mock AWS config
	client := &aws.Config{}

	// Initialize logger
	logger.InitializeLogger("info", "text", os.Stdout)

	cases := map[string]struct {
		input    string
		testCase func(*testing.T, interface{}, error)
	}{
		"AWS provider with supported service": {
			input: "ebs",
			testCase: func(t *testing.T, output interface{}, err error) {
				assert.NotNil(t, output, "service is not supported")
				assert.Implements(t, (*Cleanable)(nil), output, "returned service is not cleanable, thus it's not supported")
				assert.Nil(t, err, "error is not nil")
			},
		},
		"AWS provider with unsupported service": {
			input: "eks",
			testCase: func(t *testing.T, output interface{}, err error) {
				assert.Nil(t, output)
				if assert.Error(t, err) {
					assert.Equal(t, "service eks is not supported", err.Error())
				}
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			service, err := provider.loadService(ctx, client, test.input)
			test.testCase(t, service, err)
		})
	}
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
