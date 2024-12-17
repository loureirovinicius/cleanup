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

func TestInitialize(t *testing.T) {
	// Mock dependencies
	ctx := context.Background()

	// Mock provider
	provider := AWS{}

	// Initialize logger
	logger.InitializeLogger("info", "text", os.Stdout)

	cases := map[string]struct {
		input    string
		helpers  func()
		testCase func(*testing.T, interface{}, error)
	}{
		"Successful Initialization": {
			input: "ebs",
			helpers: func() {
				viper.Set("aws.region", "us-east-1")
			},
			testCase: func(t *testing.T, output interface{}, err error) {
				assert.Nil(t, err, "error should be nil")
			},
		},
		"Failed Initialization. Missing AWS region when loading configs": {
			input: "ebs",
			helpers: func() {
				viper.Reset()
			},
			testCase: func(t *testing.T, output interface{}, err error) {
				assert.NotNil(t, err, "error should not be nil since the AWS region is empty")
				if assert.Error(t, err, "error is nil") {
					assert.Equal(t, "AWS region can't be empty", err.Error())
				}
			},
		},
		"Failed Initialization. Application doesn't support the service requested": {
			input: "eks",
			helpers: func() {
				viper.Set("aws.region", "us-east-1")
			},
			testCase: func(t *testing.T, output interface{}, err error) {
				assert.NotNil(t, err, "error shoould not be nil since application doesn't support this AWS service")
				if assert.Error(t, err, "error is nil") {
					assert.Equal(t, "service eks is not supported", err.Error())
				}
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			test.helpers()

			err := provider.Initialize(ctx, test.input)
			test.testCase(t, nil, err)
		})
	}
}

func TestLoadAWSService(t *testing.T) {
	// Mock dependencies
	ctx := context.Background()

	// Mock provider
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

func TestCreateClient(t *testing.T) {
	// Mock dependencies
	ctx := context.Background()

	// Mock provider
	provider := AWS{}

	cases := map[string]struct {
		testCase func(*testing.T, interface{}, error)
	}{
		"AWS client created successfully": {
			testCase: func(t *testing.T, output interface{}, err error) {
				assert.Nil(t, err, "expected no error to be returned")
				assert.NotNil(t, output, "expected AWS config object to be returned")
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			client, err := provider.createClient(ctx)
			test.testCase(t, client, err)
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Mock provider
	provider := AWS{}

	cases := map[string]struct {
		helpers  func()
		testCase func(*testing.T, interface{}, error)
	}{
		"Missing AWS Region": {
			helpers: func() {
				viper.Reset()
			},
			testCase: func(t *testing.T, output interface{}, err error) {
				assert.NotNil(t, err, "expected error to be returned when region is empty")
				if assert.Error(t, err, "error is nil") {
					assert.Equal(t, "AWS region can't be empty", err.Error(), "error message should be 'AWS region can't be empty'")
				}
			},
		},
		"Valid AWS configuration": {
			helpers: func() {
				viper.Set("aws.region", "us-east-1")
			},
			testCase: func(t *testing.T, output interface{}, err error) {
				assert.Nil(t, err, "expected no error to be returned from this function")
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			test.helpers()

			err := provider.loadConfig()
			test.testCase(t, nil, err)
		})
	}
}
