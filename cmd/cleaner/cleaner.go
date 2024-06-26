package cleaner

import (
	"context"
	"fmt"
	"strings"

	"github.com/loureirovinicius/cleanup/helpers/logger"
)

func Run(ctx context.Context, services map[string]Cleanable, args []string) {
	var cmd = args[0]
	var svcName = args[1]

	service := services[svcName]

	switch cmd {
	case "list":
		res := service.List(ctx)
		val := strings.Join(res, ", ")
		logger.Log(ctx, "info", fmt.Sprintf("resources found in account: %s", val))
	case "validate":
		resources := service.List(ctx)
		for _, resource := range resources {
			empty := service.Validate(ctx, resource)
			if empty {
				msg := fmt.Sprintf("'%v' empty, can be excluded", resource)
				logger.Log(ctx, "info", msg)
			} else {
				msg := fmt.Sprintf("'%v' not empty, cannot be excluded", resource)
				logger.Log(ctx, "info", msg)
			}
		}
	case "delete":
		resources := service.List(ctx)
		for _, resource := range resources {
			empty := service.Validate(ctx, resource)
			if empty {
				res := service.Delete(ctx, resource)
				logger.Log(ctx, "info", res)
				continue
			}
			msg := fmt.Sprintf("'%v' not empty, could not delete it", resource)
			logger.Log(ctx, "info", msg)
		}
	}
}
