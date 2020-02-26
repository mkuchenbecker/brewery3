package integration

import (
	sysfsGPIO "github.com/brian-armstrong/gpio"
	"github.com/mkuchenbecker/brewery3/brewery/gpio"
	"github.com/mkuchenbecker/brewery3/brewery/utils"
	rpio "github.com/stianeikeland/go-rpio"
)

type DefaultPins struct {
}

func (d *DefaultPins) Open() error {
	return rpio.Open()
}

func (d *DefaultPins) Close() error {
	return rpio.Close()
}

func (d *DefaultPins) Pin(pin uint8) gpio.Pin {
	return rpio.Pin(pin)
}

type SysfsPins struct {
}

func (d *SysfsPins) Open() error {
	return nil
}

func (d *SysfsPins) Close() error {
	return nil
}

func (d *SysfsPins) Pin(pin uint) gpio.Pin {
	utils.Printf("getting pin: %d", pin)
	return &sysfsPin{pin: sysfsGPIO.NewOutput(pin, false)}
}

type sysfsPin struct {
	pin sysfsGPIO.Pin
}

func (pin *sysfsPin) High() {
	err := pin.pin.High()
	if err != nil {
		utils.Printf("error setting pin high: %+v", err)
	}
}

func (pin *sysfsPin) Low() {
	err := pin.pin.Low()
	if err != nil {
		utils.Printf("error setting pin low: %+v", err)
	}
}

func (pin *sysfsPin) Output() {
}
