package utils

import "fmt"

func DeferErrReturn(f func() error, err *error) {
	fnErr := f()
	if fnErr != nil {
		if *err != nil {
			*err = fmt.Errorf("recieved multiple errors: '%v' '%v'", *err, fnErr)
			return
		}
		*err = fnErr
	}
	return
}
