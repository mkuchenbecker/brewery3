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
	*err = errors.New(fmt.Sprintf("encountered a panic: %+v", r))
}
