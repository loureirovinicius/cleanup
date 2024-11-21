package providers

import (
	"context"
	"os"
	"testing"

	"github.com/loureirovinicius/cleanup/helpers/logger"
)

func TestLoadProvider(t *testing.T) {
	// Mock dependencies
	ctx := context.Background()

	// Initialize logger
	logger.InitializeLogger("info", "text", os.Stdout)

	// Test case: Unsupported provider
	t.Run("Unsupported Provider", func(t *testing.T) {
		_, err := LoadProvider(ctx, "gcp", "ebs")
		if err == nil || err.Error() != "provider gcp is not supported" {
			t.Errorf("Expected unsupported provider error, got: %v", err)
		}
	})
}
