package cleaner

import (
	"context"
	"fmt"
	"strings"

	"github.com/loureirovinicius/cleanup/helpers/logger"
	"github.com/loureirovinicius/cleanup/provider"
)

func list(ctx context.Context, service provider.Cleanable) error {
	// List all created resources for a service
	res, err := service.List(ctx)
	if err != nil {
		return fmt.Errorf("error listing resources: %v", err)
	}

	val := strings.Join(res, ", ")
	logger.Log(ctx, "info", val)

	return nil
}
