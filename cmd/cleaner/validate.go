package cleaner

import (
	"context"
	"fmt"

	"github.com/loureirovinicius/cleanup/helpers/logger"
	"github.com/loureirovinicius/cleanup/providers"
)

// Validate instances of the service passed as parameter to check whether it's being used or not
func validate(ctx context.Context, service providers.Cleanable, serviceName string) error {
	logger.Log(ctx, "info", fmt.Sprintf("Validating resources for service: %s", serviceName))

	// List all resources for the given service
	resources, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("error listing resources for service '%s': %w", serviceName, err)
	}

	// Iterate through each resource and validate
	for _, resource := range resources {
		empty, err := service.Validate(ctx, resource)
		if err != nil {
			return fmt.Errorf("error validating resource '%v' in service '%s': %w", resource, serviceName, err)
		}

		// Log whether each resource can be excluded based on validation
		if empty {
			logger.Log(ctx, "info", fmt.Sprintf("Resource '%v' in service '%s' is empty and can be excluded.", resource, serviceName))
		} else {
			logger.Log(ctx, "info", fmt.Sprintf("Resource '%v' in service '%s' is not empty and cannot be excluded.", resource, serviceName))
		}
	}

	logger.Log(ctx, "debug", fmt.Sprintf("Validation completed for service: %s", serviceName))
	return nil
}
