package providers

import (
	"context"
	"os"
	"testing"

	"github.com/loureirovinicius/cleanup/helpers/logger"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadProvider(t *testing.T) {
	// Mock dependencies
	ctx := context.Background()

	// Initialize logger
	logger.InitializeLogger("info", "text", os.Stdout)

	cases := map[string]struct {
		input    string
		helpers  func()
		testCase func(*testing.T, interface{}, error)
	}{
		"Supported Provider (AWS)": {
			input: "aws",
			helpers: func() {
				viper.Set("aws.region", "us-east-1")
			},
			testCase: func(t *testing.T, output interface{}, err error) {
				assert.NotNil(t, output, "provider is not supported")
				assert.Nil(t, err, "error is not nil")
			},
		},
		"Unsupported Provider": {
			input:   "azure",
			helpers: func() {},
			testCase: func(t *testing.T, output interface{}, err error) {
				assert.Nil(t, output)
				if assert.Error(t, err) {
					assert.Equal(t, "provider azure is not supported", err.Error())
				}
			},
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			test.helpers()

			service, err := LoadProvider(ctx, test.input, "ebs")
			test.testCase(t, service, err)
		})
	}
}
