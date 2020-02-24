package gpio

//go:generate mockgen -destination=./mocks/mock_gpio.go github.com/mkuchenbecker/brewery3/brewery/gpio IGPIO,Pin,Temperature

// TemperatureAddress is a wrapper around string for stricter typing.
type TemperatureAddress string

// IGPIO is the interface for interacting with the github.com/stianeikeland/go-rpio library.
type IGPIO interface {
	Open() error
	Close() error
	Pin(uint8) Pin
}

// Pin is the interface for interacting with a pin from github.com/stianeikeland/go-rpio.
type Pin interface { //nolint:golint
	High()
	Low()
	Output()
}

// Temperature is the interface for interacting with github.com/yryz/ds18b20
type Temperature interface { //nolint:golint
	Sensors() ([]string, error)
	Temperature(TemperatureAddress) (float64, error)
}
