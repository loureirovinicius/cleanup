package cleaner

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/loureirovinicius/cleanup/helpers/logger"
)

type CleanerListCommand struct{}

func (c *CleanerListCommand) Run(ctx context.Context, service Cleanable) error {
	var value string

	// Creates flag parameters for the "List" operation
	list := flag.NewFlagSet("list", flag.ExitOnError)
	list.StringVar(&value, "service", "", "cloud provider service")
	list.StringVar(&value, "s", "", "cloud provider service")
	err := list.Parse(os.Args[2:])
	if err != nil {
		return fmt.Errorf("error parsing CLI args: %v", err)
	}

	// List all created resources for a service
	res, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("error listing resources: %v", err)
	}

	val := strings.Join(res, ", ")
	logger.Log(ctx, "info", val)

	return nil
}
