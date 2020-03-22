package cli

import (
	"fmt"
	"time"

	"github.com/mkuchenbecker/brewery3/brewery/utils"

	"github.com/pkg/errors"

	"context"
	"strconv"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"

	"github.com/urfave/cli"
)

const (
	// MaxBoil is the maximum boil kettle temerature.
	MaxBoil = 95
	// HERMSTolerance is temperature above the mash target the HERMS temperature is allowed to drift above. A higher herms tolerance will speed heating at risk of overshooting the temperature.
	HERMSTolerance = 3
	// MashTolerance is the maximum drift above the target for the mash kettle before turning off the element. Note: a HERMSTolerance setto high will still cause the mash to overshoot.
	MashTolerance = 0.5
)

func parseTemp(in string) (float64, error) {
	t, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("unable to parse temperature '%s'", in))
	}
	return t, nil
}

func getMashRequest(temp float64) *model.ControlRequest {
	utils.Print(fmt.Sprintf("Mashing @ %fC", temp))
	return &model.ControlRequest{Scheme: &model.ControlScheme{
		Scheme: &model.ControlScheme_Mash_{
			Mash: &model.ControlScheme_Mash{
				HermsMaxTemp: temp + HERMSTolerance,
				HermsMinTemp: temp,
				MashMinTemp:  temp,
				MashMaxTemp:  temp + MashTolerance,
				BoilMinTemp:  temp,
				BoilMaxTemp:  MaxBoil,
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
	utils.Print("Starting CLI\n")
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "mash",
			Value: "float64",
			Usage: "controls the temperature of the mash tun",
		},
		&cli.BoolFlag{
			Name:  "boil",
			Usage: "turn the element on",
		},
		&cli.BoolFlag{
			Name:  "off",
			Usage: "turns the element off",
		},
	}

	app.Name = "Brewery3 CLI"
	app.Usage = "Control the brewery with a CLI"
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
