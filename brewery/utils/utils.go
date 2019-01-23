package utils

import (
	"fmt"
	"time"
)

const UpdateInterval = 5 * time.Second

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

type Printer interface {
	Print(s string)
}

type DefualtPrinter struct {
}

func (p *DefualtPrinter) Print(s string) {
	Print(s)
}

func Print(s string) {
	fmt.Printf("%s - %s\n", time.Now().Format(time.StampMilli), s)
}

func LogError(p Printer, err error, msg string) {
	var printer Printer = &DefualtPrinter{}
	if p != nil {
		printer = p
	}
	printer.Print(fmt.Sprintf("%s : %s", msg, err.Error()))
}

func BackgroundErrReturn(f func() error) {
	err := f()
	if err != nil {
		LogError(nil, err, "background function encountered error")
	}
}
