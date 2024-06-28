package cleaner

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/loureirovinicius/cleanup/helpers/logger"
)

type CleanerValidateCommand struct{}

func (c *CleanerValidateCommand) Run(ctx context.Context, service Cleanable) error {
	var value string

	// Creates flag parameters for the "Validate" operation
	validate := flag.NewFlagSet("validate", flag.ExitOnError)
	validate.StringVar(&value, "service", "", "cloud provider service")
	validate.StringVar(&value, "s", "", "cloud provider service")
	validate.Parse(os.Args[2:])

	resources, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("error listing resources: %v", err)
	}

	for _, resource := range resources {
		empty, err := service.Validate(ctx, resource)
		if err != nil {
			return fmt.Errorf("error validating resources: %v", err)
		}

		if empty {
			msg := fmt.Sprintf("'%v' empty, can be excluded", resource)
			logger.Log(ctx, "info", msg)
		} else {
			msg := fmt.Sprintf("'%v' not empty, cannot be excluded", resource)
			logger.Log(ctx, "info", msg)
		}
	}

	return nil
}
