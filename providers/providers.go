package providers

import (
	"context"
	"fmt"
)

// Initialize the cloud provider being used during the execution
func LoadProvider(ctx context.Context, provider string, service string) (Cleanable, error) {
	switch provider {
	case "aws":
		// Instantiate AWS provider
		aws := AWS{}

		// Initialize the provider with required configs and specified service
		if err := aws.Initialize(ctx, service); err != nil {
			return nil, fmt.Errorf("error initializing AWS functions. Reason: %v", err)
		}

		// Return only the requested service in the map
		return aws.Service, nil
	default:
		return nil, fmt.Errorf("provider %s is not supported", provider)
	}
}
