package cleaner

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/loureirovinicius/cleanup/helpers/logger"
)

type CleanerDeleteCommand struct{}

func (c *CleanerDeleteCommand) Run(ctx context.Context, service Cleanable) error {
	var value string

	// Creates flag parameters for the "Validate" operation
	validate := flag.NewFlagSet("validate", flag.ExitOnError)
	validate.StringVar(&value, "service", "", "cloud provider service")
	validate.StringVar(&value, "s", "", "cloud provider service")
	validate.Parse(os.Args[2:])

	// List all created resources for a service
	resources, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("error listing resources: %v", err)
	}

	// Loop through all the resources validating them
	for _, resource := range resources {
		empty, err := service.Validate(ctx, resource)
		if err != nil {
			return err
		}

		if empty {
			msg := fmt.Sprintf("'%v' empty, can be excluded", resource)
			logger.Log(ctx, "info", msg)
		} else {
			msg := fmt.Sprintf("'%v' not empty, cannot be excluded", resource)
			logger.Log(ctx, "info", msg)
			continue
		}

		// Delete the resource
		err = service.Delete(ctx, resource)
		if err != nil {
			return fmt.Errorf("error deleting resource: %v", err)
		}
	}

	return nil
}
