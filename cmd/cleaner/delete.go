package cleaner

import (
	"context"
	"fmt"

	"github.com/loureirovinicius/cleanup/helpers/logger"
	"github.com/loureirovinicius/cleanup/providers"
)

// Delete unused instances of the service passed as parameter
func delete(ctx context.Context, service providers.Cleanable, serviceName string) error {
	logger.Log(ctx, "info", fmt.Sprintf("Deleting resources for service: %s", serviceName))

	// List all resources for the given service
	resources, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("error listing resources for service '%s': %w", serviceName, err)
	}

	// Iterate through each resource to validate and delete if empty
	for _, resource := range resources {
		empty, err := service.Validate(ctx, resource)
		if err != nil {
			return fmt.Errorf("error validating resource '%v' in service '%s': %w", resource, serviceName, err)
		}

		if empty {
			logger.Log(ctx, "info", fmt.Sprintf("Resource '%v' in service '%s' is empty and can be excluded.", resource, serviceName))

			// Attempt to delete the empty resource
			err = service.Delete(ctx, resource)
			if err != nil {
				return fmt.Errorf("error deleting resource '%v' in service '%s': %w", resource, serviceName, err)
			}
			logger.Log(ctx, "info", fmt.Sprintf("Resource '%v' in service '%s' has been deleted successfully.", resource, serviceName))
		} else {
			logger.Log(ctx, "info", fmt.Sprintf("Resource '%v' in service '%s' is not empty and cannot be excluded.", resource, serviceName))
		}
	}

	logger.Log(ctx, "info", fmt.Sprintf("Deletion completed for service: %s", serviceName))
	return nil
}
