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
			*err = fmt.Errorf("received multiple errors: '%v' '%v'", *err, fnErr)
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

func Printf(format string, i ...interface{}) {
	Print(fmt.Sprintf(format, i...))
}

func LogError(p Printer, err error, msg string) {
	var printer Printer = &DefualtPrinter{}
	if p != nil {
		printer = p
	}
	printer.Print(fmt.Sprintf("%s : %s", msg, err.Error()))
}

func LogIfError(p Printer, err error, msg string) bool {
	if err == nil {
		return false
	}
	LogError(p, err, msg)
	return true
}

func BackgroundErrReturn(p Printer, f func() error) {
	err := f()
	if err != nil {
		LogError(p, err, "background function encountered error")
	}
}

func RunLoop(ttl time.Duration, minLoopInterval time.Duration, fn func() error) (err error) {
	start := time.Now()
	for {
		now := time.Now()
		if now.Sub(start) > ttl {
			return
		}
		Print("[RunLoop]")
		err := fn()
		if err != nil {
			Print(fmt.Sprintf("[RunLoop] %s", err.Error()))
		}
		sleepTime := minLoopInterval - time.Since(now)
		if sleepTime > time.Duration(0) {
			time.Sleep(sleepTime)
		}
		err = nil
	}
}

type TemperaturePointSink struct {
	Temps map[string]float64
}

func NewTemperaturePointSink() *TemperaturePointSink {
	return &TemperaturePointSink{
		Temps: make(map[string]float64),
	}
}

func (tps *TemperaturePointSink) ToInterface() map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range tps.Temps {
		res[k] = interface{}(v)
	}
	return res
}

func (tps *TemperaturePointSink) Log() {
	for k, temp := range tps.Temps {
		Printf("Temp [%s]%f", k, temp)
	}
}
