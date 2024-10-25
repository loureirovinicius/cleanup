package cleaner

import (
	"context"
	"fmt"

	"github.com/loureirovinicius/cleanup/helpers/logger"
	"github.com/loureirovinicius/cleanup/provider"
)

func validate(ctx context.Context, service provider.Cleanable) error {
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
