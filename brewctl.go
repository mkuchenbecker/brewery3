package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/utils"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

func MakeBreweryClient(port int) model.BreweryClient {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	//defer conn.Close() //TODO
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
	app := cli.NewApp()

	client := MakeBreweryClient(8100)

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "mash",
			Value: "float64",
			Usage: "temperature",
		}}

	app.Name = "mash"
	app.Usage = "mash the beer!"
	app.Action = func(c *cli.Context) error {
		if !c.IsSet("mash") {
			utils.Print(fmt.Sprintf("Must set temp flag!"))
			os.Exit(1)
		}
		temp, err := parseTemp(c.String("mash"))
		if err != nil {
			utils.Print(fmt.Sprintf("%+v", err))
			os.Exit(1)
		}

		utils.Print(fmt.Sprintf("Input Temp: %s", c.String("temp")))
		utils.Print(fmt.Sprintf("Mashing @ %fC", temp))

		_, err = client.Control(context.Background(),
			&model.ControlRequest{Scheme: &model.ControlScheme{
				Scheme: &model.ControlScheme_Mash_{
					Mash: &model.ControlScheme_Mash{
						HermsMaxTemp: temp + 15,
						HermsMinTemp: temp,
						MashMinTemp:  temp,
						MashMaxTemp:  temp + .5,
						BoilMinTemp:  temp,
						BoilMaxTemp:  100,
					},
				},
			},
			})
		return err
	}

	err := app.Run(os.Args)
	if err != nil {
		utils.Print("encountered an error:")
		log.Fatal(err)
	}
}
