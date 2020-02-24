package cli

import (
	"fmt"
	"time"

	"github.com/mkuchenbecker/brewery3/brewery/utils"

	"context"
	"strconv"
	"strings"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"

	"github.com/urfave/cli"
)

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

func getMashRequest(temp float64) *model.ControlRequest {
	utils.Print(fmt.Sprintf("Mashing @ %fC", temp))

	return &model.ControlRequest{Scheme: &model.ControlScheme{
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
	}
}

func getOffRequest() *model.ControlRequest {
	return &model.ControlRequest{
		Scheme: &model.ControlScheme{
			Scheme: &model.ControlScheme_Off_{
				Off: &model.ControlScheme_Off{},
			},
		},
	}
}

func getBoilRequest() *model.ControlRequest {
	return &model.ControlRequest{
		Scheme: &model.ControlScheme{
			Scheme: &model.ControlScheme_Boil_{
				Boil: &model.ControlScheme_Boil{},
			},
		},
	}
}

// Run takes a command to det the empterature of the mash or boil server.
func Run(client model.BreweryClient, args []string) error {
	utils.Print("received CLI command\n")
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "mash",
			Value: "float64",
			Usage: "temperature",
		},
		&cli.BoolFlag{
			Name:  "boil",
			Usage: "temperature",
		},
		&cli.BoolFlag{
			Name:  "off",
			Usage: "off",
		},
	}

	app.Name = "brew"
	app.Usage = "brew the beer!"
	app.Action = func(c *cli.Context) error {
		var req *model.ControlRequest
		var err error
		if c.IsSet("mash") {
			utils.Print("mash\n")
			var temp float64
			if temp, err = parseTemp(c.String("mash")); err != nil || temp < 0 || temp > 100 {
				err = fmt.Errorf("mash temp invalid, must be a float from 0-100: %s", c.String("mash"))
				utils.Print(err.Error())
				return err
			}

			utils.Printf("mash: %f\n", temp)
			req = getMashRequest(temp)
		}
		if c.IsSet("boil") {
			utils.Print("boil\n")
			req = getBoilRequest()
		}
		if c.IsSet("off") {
			utils.Print("boil\n")
			req = getOffRequest()
		}
		if req == nil {
			return fmt.Errorf("no arguments specified")
		}

		clientDeadline := time.Now().Add(time.Duration(60000) * time.Millisecond)
		ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
		defer cancel()
		utils.Printf("sending request:\n%+v\n", req)
		_, err = client.Control(ctx, req)
		return err
	}
	return app.Run(args)
}
