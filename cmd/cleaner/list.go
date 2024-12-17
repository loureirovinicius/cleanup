package cleaner

import (
	"context"
	"fmt"
	"strings"

	"github.com/loureirovinicius/cleanup/helpers/logger"
	"github.com/loureirovinicius/cleanup/providers"
)

// List instances of the service passed as parameter
func list(ctx context.Context, service providers.Cleanable, serviceName string) error {
	// List all created resources for a service
	logger.Log(ctx, "info", fmt.Sprintf("Listing resources for service: %s", serviceName))

	resources, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("error listing resources for service '%s': %w", serviceName, err)
	}

	// Join resource IDs or names into a single string for logging
	resourceList := strings.Join(resources, ", ")
	logger.Log(ctx, "info", fmt.Sprintf("Resources for %s: %s", serviceName, resourceList))

	return nil
}
