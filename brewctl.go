package main

import (
	"context"
	"fmt"
	"time"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/rpi/sensors"
	"github.com/mkuchenbecker/brewery3/brewery/utils"
	"google.golang.org/grpc"
)

func MakeTemperatureClient(port int, address string) model.ThermometerClient {
	utils.Print(fmt.Sprintf("Starting temperature server on port: %d", port))
	go sensors.StartThermometer(port, address)
	utils.Print(fmt.Sprintf("Waiting for discovery on port: %d", port))
	time.Sleep(5 * time.Second)
	utils.Print(fmt.Sprintf("Connecting to client: %d", port))
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := model.NewThermometerClient(conn)
	res, err := client.Get(context.Background(), &model.GetRequest{})
	if err != nil {
		panic(err)
	}
	utils.Print(fmt.Sprintf("temp: %f", res.Temperature))
	return client
}

// func MakeSwitchClient(port int, pin uint8) model.SwitchClient {
// 	utils.Print(fmt.Sprintf("Starting switch server on port: %d", port)
// 	go element.StartHeater(port, pin)
// 	utils.Print(fmt.Sprintf("Waiting for discovery on port: %d", port)
// 	time.Sleep(5 * time.Second)
// 	utils.Print(fmt.Sprintf("Connecting to client: %d", port)
// 	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer conn.Close()
// 	client := model.NewSwitchClient(conn)
// 	_, err = client.Off(context.Background(), &model.OffRequest{})
// 	if err != nil {
// 		panic(err)
// 	}
// 	return client
// }

// func MakeBreweryClient(port int) model.BreweryClient {
// 	conn, err := grpc.Dial(fmt.Sprintf("localhost://%d", port))
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer conn.Close()
// 	return model.NewBreweryClient(conn)
// }

// func parseTemp(in string) (float64, error) {
// 	i := strings.Index(in, "f")
// 	str := in
// 	if i > -1 {
// 		str = str[:i]
// 	}
// 	t, err := strconv.ParseFloat(str, 64)
// 	if err != nil {
// 		return 0, err
// 	}

// 	if i > -1 {
// 		return (t - 32) * 5 / 9, nil
// 	}

// 	return t, nil
// }

const (
	MashAddr  = "28-0315712c08ff"
	HermsAddr = "28-0315715039ff"
	BoilAddr  = "28-031571188aff"
)

func main() {
	MakeTemperatureClient(8090, MashAddr)
	MakeTemperatureClient(8091, HermsAddr)
	MakeTemperatureClient(8092, BoilAddr)

	// brewery := rpi.Brewery{
	// 	MashSensor:  MakeTemperatureClient(8090, "28-0315715039ff"),
	// 	HermsSensor: MakeTemperatureClient(8091, "28-0315712c08ff"),
	// 	BoilSensor:  MakeTemperatureClient(8092, "28-031571188aff"),
	// 	Element:     MakeSwitchClient(8110, 11),
	// }
	// rpi.StartBrewery(8100, &brewery)

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
	// 		utils.Print(fmt.Sprintf("Must set temp flag!")
	// 		os.Exit(1)
	// 	}
	// 	temp, err := parseTemp(c.String("mash"))
	// 	if err != nil {
	// 		utils.Print(fmt.Sprintf("%+v", err)
	// 		os.Exit(1)
	// 	}

	// 	utils.Print(fmt.Sprintf("Input Temp: %s", c.String("temp"))
	// 	utils.Print(fmt.Sprintf("Mashing @ %fC", temp)

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
