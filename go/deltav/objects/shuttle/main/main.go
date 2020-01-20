package main

import (
	"context"
	"log"

	protos "github.com/golang001/deltav/model/gomodel"
	"github.com/golang001/deltav/objects/shuttle"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := protos.NewWorldModelClient(conn)

	vessel := protos.Vessel{
		Sensors: &protos.SensorSystem{
			Sensors: []*protos.Sensor{
				{
					SensorType:    protos.Sensor_PASSIVE,
					RadiationType: protos.RadiationType_GAMMA,
				},
			},
		},
	}
	_, err = shuttle.NewPlayerShuttle(vessel, c)
	if err != nil {
		log.Fatalf("could not create user: %v", err)

	}
	r, err := c.Register(context.Background(), &protos.RegisterRequest{})
	if err != nil {
		log.Fatalf("could not register: %v", err)
	}
	log.Printf("Greeting: %s", r.Effect)
}
