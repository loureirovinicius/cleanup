package cleaner

import (
	"context"
	"errors"
	"fmt"
	"os"
)

func Run(ctx context.Context, services map[string]Cleanable) error {
	var cmd, svcName string

	// Checks if it contains the minimum required args
	if len(os.Args) < 3 {
		usage()
		return nil
	}

	cmd, svcName = os.Args[1], os.Args[3]

	// Validates if service exists so that the program doesn't crash
	if _, ok := services[svcName]; !ok {
		return errors.New("please provide an existing service. Refer to this project's README for the full list of available services")
	}
	service := services[svcName]

	switch cmd {
	case "list":
		return (&CleanerListCommand{}).Run(ctx, service)
	case "validate":
		return (&CleanerValidateCommand{}).Run(ctx, service)
	case "delete":
		return (&CleanerDeleteCommand{}).Run(ctx, service)
	default:
		usage()
		return nil
	}
}

func usage() {
	fmt.Print(`Cleanup is a tool for cloud providers' resources cleaning. You can quickly list, validate or delete resources from your current cloud provider vendor. The tool is being incrementally built, so the only provider currently supported is AWS with few resources.

Usage:
	cleanup <command> (--service | -service) <service_name>

Commands Usage:
	cleanup list (--service | -service) <service_name>

		Options:
			--service STRING  (required) the service name you're trying to list (these service names are available in the documentation).

	cleanup validate (--service | -service) <service_name>

		Options:
			--service STRING  (required) the service name you're trying to validate (these service names are available in the documentation). Each resource has its own rules to be considered empty, so check docs to understand these rules.

	cleanup delete (--service | -service) <service_name>

		Options:
			--service STRING  (required) the service name you're trying to delete (these service names are available in the documentation). A resource can only be deleted if empty, so check it first using the "validate" operation.
	`)
}
