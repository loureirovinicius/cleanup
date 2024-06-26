package cleaner

import (
	"context"
)

type Cleanable interface {
	List(context.Context) []string
	Validate(context.Context, string) bool
	Delete(context.Context, string) string
}
