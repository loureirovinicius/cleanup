package provider

import (
	"context"

	"github.com/loureirovinicius/cleanup/cmd/cleaner"
)

type ProviderConfig struct {
	AWS
	// GCP (soon)
}

type Provider interface {
	Initialize(context.Context, *ProviderConfig)
}

func LoadProvider(ctx context.Context, provider string) map[string]cleaner.Cleanable {
	cfg := ProviderConfig{}

	switch provider {
	case "aws":
		cfg.AWS.Initialize(ctx, &cfg)
		p := cfg.AWS
		return p.Resources
	}

	return nil
}
