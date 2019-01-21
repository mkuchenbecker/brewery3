package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/rpi"
	"github.com/mkuchenbecker/brewery3/brewery/rpi/element"
	"github.com/mkuchenbecker/brewery3/brewery/rpi/sensors"
	"google.golang.org/grpc"
)

func MakeTemperatureClient(port int, address string) model.ThermometerClient {
	fmt.Printf("Starting temperature server on port: %d/n", port)
	go sensors.StartThermometer(port, address)
	fmt.Printf("Waiting for discovery on port: %d/n", port)
	time.Sleep(5 * time.Second)
	fmt.Printf("Connecting to client: %d/n", port)
	conn, err := grpc.Dial(fmt.Sprintf("localhost://%d", port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := model.NewThermometerClient(conn)
	_, err = client.Get(context.Background(), &model.GetRequest{})
	if err != nil {
		panic(err)
	}
	return client
}

func MakeSwitchClient(port int, pin int) model.SwitchClient {
	fmt.Printf("Starting switch server on port: %d/n", port)
	go element.StartHeater(port, pin)
	fmt.Printf("Waiting for discovery on port: %d/n", port)
	time.Sleep(5 * time.Second)
	fmt.Printf("Connecting to client: %d/n", port)
	conn, err := grpc.Dial(fmt.Sprintf("localhost://%d", port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := model.NewSwitchClient(conn)
	_, err = client.Off(context.Background(), &model.OffRequest{})
	if err != nil {
		panic(err)
	}
	return client
}

func MakeBreweryClient(port int) model.BreweryClient {
	conn, err := grpc.Dial(fmt.Sprintf("localhost://%d", port))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	return model.NewBreweryClient(conn)
}

func parseTemp(in string) (float64, error) {
	i := strings.Index(in, "f")
	str := in
	if i > -1 {
		str = str[:i]
	}
	t, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}

	if i > -1 {
		return (t - 32) * 5 / 9, nil
	}

	return t, nil
}

func main() {

	brewery := rpi.Brewery{
		MashSensor:  MakeTemperatureClient(8090, "28-031571188aff"),
		HermsSensor: MakeTemperatureClient(8091, "28-0315712c08ff"),
		BoilSensor:  MakeTemperatureClient(8092, "28-0315715039ff"),
		Element:     MakeSwitchClient(8110, 11),
	}
	rpi.StartBrewery(8100, &brewery)

	// app := cli.NewApp()

	// client := MakeBreweryClient(8100)

	// app.Flags = []cli.Flag{
	// 	cli.StringFlag{
	// 		Name:  "mash",
	// 		Value: "float64",
	// 		Usage: "temperature",
	// 	}}

	// app.Name = "mash"
	// app.Usage = "mash the beer!"
	// app.Action = func(c *cli.Context) error {
	// 	if !c.IsSet("mash") {
	// 		fmt.Printf("Must set temp flag!\n")
	// 		os.Exit(1)
	// 	}
	// 	temp, err := parseTemp(c.String("mash"))
	// 	if err != nil {
	// 		fmt.Printf("%+v", err)
	// 		os.Exit(1)
	// 	}

	// 	fmt.Printf("Input Temp: %s\n", c.String("temp"))
	// 	fmt.Printf("Mashing @ %fC\n", temp)

	// 	_, err = client.Control(context.Background(),
	// 		&model.ControlRequest{Scheme: &model.ControlScheme{
	// 			Scheme: &model.ControlScheme_Mash_{
	// 				Mash: &model.ControlScheme_Mash{
	// 					HermsMaxTemp: temp + 15,
	// 					HermsMinTemp: temp,
	// 					MashMinTemp:  temp,
	// 					MashMaxTemp:  temp + .5,
	// 					BoilMinTemp:  temp,
	// 					BoilMaxTemp:  100,
	// 				},
	// 			},
	// 		},
	// 		})
	// 	return err
	// }

	// err := app.Run(os.Args)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
