package utils

import (
	"fmt"
	"time"
)

func DeferErrReturn(f func() error, err *error) {
	fnErr := f()
	if fnErr != nil {
		if *err != nil {
			*err = fmt.Errorf("recieved multiple errors: '%v' '%v'", *err, fnErr)
			return
		}
		*err = fnErr
	}
}

func Print(s string) {
	fmt.Printf("%s - %s\n", time.Now().Format(time.StampMilli), s)
}
