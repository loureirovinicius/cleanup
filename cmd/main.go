package main

import (
	"context"

	"github.com/loureirovinicius/cleanup/cmd/cleaner"
	"github.com/loureirovinicius/cleanup/cmd/cli"
	"github.com/loureirovinicius/cleanup/config"
	"github.com/loureirovinicius/cleanup/helpers/logger"
	"github.com/loureirovinicius/cleanup/provider"
)

func main() {
	ctx := context.Background()

	err := config.Start()
	if err != nil {
		logger.Log(ctx, "error", err.Error())
	}
	args := cli.Start()

	// Load the provider configuration
	services, err := provider.LoadProvider(ctx, "aws")
	if err != nil {
		logger.Log(ctx, "error", err.Error())
	}

	// Starts the Cleaner
	cleaner.Run(ctx, services, args)
}
