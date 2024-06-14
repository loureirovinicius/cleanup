package cleaner

import (
	"context"
)

func Run(ctx context.Context, services map[string]Cleanable, args []string) []string {
	var cmd = args[0]
	var svcName = args[1]

	service := services[svcName]

	switch cmd {
	case "list":
		res := service.List(ctx)
		return res
	// case "validate":
	// 	resources := service.List(ctx)
	// 	for _, resource := range resources {
	// 		empty := service.Validate(resource)
	// 		return empty
	// 	}
	// case "delete":
	// 	resources := service.List(ctx)
	// 	for _, resource := range resources {
	// 		empty := service.Validate()
	// 		if empty {
	// 			res := service.Delete(resource)
	// 			return res
	// 		}
	// 		return
	// 	}
	case "":
		usage()
	}

	return nil
}

func usage() {}
