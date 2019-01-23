package gpio

//go:generate mockgen -destination=./mocks/mock_gpio.go github.com/mkuchenbecker/brewery3/brewery/gpio IGPIO,Controller,GPIOPin

// TemperatureAddress is a wrapper around string for stricter typing.
type TemperatureAddress string

// Controller is an interface for interacting with the raspberry pi GPIO.
// Any functionality that interfaces with hardware should be made an interface here
// so it can be mocked out.
type Controller interface {
	PowerPin(pin uint8, on bool) error
	ReadTemperature(address TemperatureAddress) (float64, error)
}

// IGPIO is the interface for interacting with the github.com/stianeikeland/go-rpio library.
type IGPIO interface {
	Open() error
	Close() error
	Pin(uint8) GPIOPin
}

// GPIOPin is the interface for interacting with a pin from github.com/stianeikeland/go-rpio.
type GPIOPin interface { //nolint:golint
	High()
	Low()
}

// GPIOTemperature is the interface for interacting with github.com/yryz/ds18b20
type GPIOTemperature interface { //nolint:golint
	Sensors() ([]string, error)
	Temperature(TemperatureAddress) (float64, error)
}
