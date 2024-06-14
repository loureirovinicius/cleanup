package cleaner

import (
	"context"
)

type Cleanable interface {
	List(context.Context) []string
	Validate(Cleanable) bool
	Delete(string) string
}
