package utils

import (
	"context"
	"fmt"
)

// CapturePanic should be deferred and will capture any panics in the foreground process.
func CapturePanic(ctx context.Context, res interface{}, err *error) {
	r := recover()
	if r != nil {
		res = nil
		*err = fmt.Errorf("encountered a panic: %+v", r)
	}
}
