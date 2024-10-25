package provider

import (
	"context"
)

type Cleanable interface {
	List(context.Context) ([]string, error)
	Validate(context.Context, string) (bool, error)
	Delete(context.Context, string) error
}
