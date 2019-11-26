package utils

import (
	"errors"
	"fmt"
)

func PanicRecover(err *error) {
	r := recover()
	if r == nil {
		return
	}
	// errors.New grabs a stacktrace, fmt.Errorf does not. Golint error prefers the latter.
	*err = errors.New(fmt.Sprintf("encountered a panic: %+v", r)) //nolint:golint,gosimple
}
