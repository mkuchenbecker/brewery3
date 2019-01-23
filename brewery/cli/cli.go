package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/utils"
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

func getPowerRequst(power int64) *model.ControlRequest {
	return &model.ControlRequest{Scheme: &model.ControlScheme{
		Scheme: &model.ControlScheme_Power_{
			Power: &model.ControlScheme_Power{
				PowerLevel: float64(power),
			},
		},
	},
	}
}

func getBoilRequest() *model.ControlRequest {
	return getPowerRequst(75)
}

func Run(client model.BreweryClient, args []string) error {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "mash",
			Value: "float64",
			Usage: "temperature",
		},
		cli.Uint64Flag{
			Name:  "power",
			Usage: "power percent",
		},
		cli.BoolFlag{
			Name:  "boil",
			Usage: "temperature",
		},
	}

	app.Name = "brew"
	app.Usage = "brew the beer!"
	app.Action = func(c *cli.Context) error {
		var req *model.ControlRequest
		var err error
		if c.IsSet("mash") {
			var temp float64
			if temp, err = parseTemp(c.String("mash")); err != nil || temp < 0 || temp > 100 {
				err = fmt.Errorf("mash temp invalid, must be a float from 0-100: %s", c.String("mash"))
				utils.Print(err.Error())
				return err
			}
			req = getMashRequest(temp)
		}
		if c.IsSet("boil") {
			req = getBoilRequest()
		}
		if c.IsSet("power") {
			power := c.Uint64("power")
			if power > 100 || power < 0 {
				err = fmt.Errorf("power level invalid, must be an integer from 0-100: %d", power)
				utils.Print(err.Error())
				return err
			}
			req = getPowerRequst(int64(power))
		}
		if req == nil {
			return fmt.Errorf("no arguments specified")
		}

		_, err = client.Control(context.Background(), req)
		return err
	}
	return app.Run(args)
}
