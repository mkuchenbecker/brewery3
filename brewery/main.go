package main

import (
	"log"
	"time"

	"github.com/mkuchenbecker/brewery3/brewery/gpio/integration"
)

func main() {
	pins := integration.DefaultPins{}
	if err := pins.Open(); err != nil {
		log.Fatalf("failed to open gpio: %v", err)
	}
	defer pins.Close()
	pin := pins.Pin(uint8(5))
	pin.Output()
	for {
		pin.High()
		time.Sleep(time.Second)
		pin.Low()
		time.Sleep(time.Second)
	}
}
