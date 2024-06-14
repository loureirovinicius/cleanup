package main

import (
	"context"
	"fmt"

	"github.com/loureirovinicius/cleanup/cmd/cleaner"
	"github.com/loureirovinicius/cleanup/cmd/cli"
	"github.com/loureirovinicius/cleanup/config"
	"github.com/loureirovinicius/cleanup/provider"
)

func main() {
	ctx := context.TODO()

	config.Start()
	args := cli.Start()

	// Load the provider configuration
	services := provider.LoadProvider(ctx, "aws")

	// Starts the Cleaner
	res := cleaner.Run(ctx, services, args)

	fmt.Println(res)
}
