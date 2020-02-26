package main

import (
	"log"
	"os"
	"strconv"

	"github.com/mkuchenbecker/brewery3/brewery/element/element"
)

func main() { // nolint: deadcode
	strPort := os.Getenv("PORT")
	strPin := os.Getenv("GPIO_PIN")

	port, err := strconv.ParseInt(strPort, 10, 32)
	if err != nil {
		log.Fatalf("Invalid port is not 32 bit int: %s", strPort)
	}

	pinNum, err := strconv.ParseInt(strPin, 10, 8)
	if err != nil {
		log.Fatalf("invalid pin : %s", strPin)
	}

	element.StartElement(pinNum, port)

}
