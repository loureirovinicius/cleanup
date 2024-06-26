package provider

import (
	"context"
	"fmt"

	"github.com/loureirovinicius/cleanup/cmd/cleaner"
)

type ProviderConfig struct {
	AWS
	// GCP (soon)
}

type Provider interface {
	Initialize(context.Context, *ProviderConfig) error
}

func LoadProvider(ctx context.Context, provider string) (map[string]cleaner.Cleanable, error) {
	cfg := ProviderConfig{}

	switch provider {
	case "aws":
		err := cfg.AWS.Initialize(ctx, &cfg)
		if err != nil {
			return nil, fmt.Errorf("error initializing AWS functions. Reason: %v", err)
		}

		p := cfg.AWS
		return p.Resources, nil
	}

	return nil, nil
}
