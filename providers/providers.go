package providers

import (
	"context"
	"fmt"
)

type ProviderConfig struct {
	AWS
}

type Provider interface {
	Initialize(context.Context, *ProviderConfig, string) error
}

// Initialize the cloud provider being used during the execution
func LoadProvider(ctx context.Context, provider string, service string) (Cleanable, error) {
	cfg := ProviderConfig{}

	switch provider {
	case "aws":
		// Initialize the provider with required configs and specified service
		err := cfg.AWS.Initialize(ctx, &cfg, service)
		if err != nil {
			return nil, fmt.Errorf("error initializing AWS functions. Reason: %v", err)
		}

		// Return only the requested service in the map
		return cfg.AWS.Service, nil
	default:
		return nil, fmt.Errorf("provider %s is not supported", provider)
	}
}
