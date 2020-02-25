package brewery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mkuchenbecker/brewery3/brewery/logger"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/utils"
)

// Brewery is the central controller for various heating elements, sensors, and switches.
// It coordinates all actions of the brewery.
type Brewery struct {
	Scheme *model.ControlScheme
	mux    sync.RWMutex

	MashSensor  model.ThermometerClient
	BoilSensor  model.ThermometerClient
	HermsSensor model.ThermometerClient

	Logger logger.Logger

	Element model.SwitchClient
}

// Control implements the Brewery.Control method. Given a ControlScheme
// The brewery will then attempt to follow the scheme.
func (b *Brewery) Control(ctx context.Context,
	req *model.ControlRequest) (res *model.ControlResponse, err error) {
	utils.Print("Received control request")
	b.replaceConfig(req.Scheme)
	return &model.ControlResponse{}, nil
}

func (b *Brewery) replaceConfig(scheme *model.ControlScheme) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Scheme = scheme
}

func (b *Brewery) mashThermOn(ctx context.Context) (bool, error) {
	tps := utils.NewTemperaturePointSink()
	resBoil, err := b.BoilSensor.Get(ctx, &model.GetRequest{})
	if err != nil {
		return false, err
	}
	tps.Temps["boil"] = resBoil.Temperature

	resHerms, err := b.HermsSensor.Get(ctx, &model.GetRequest{})
	if err != nil {
		return false, err
	}
	tps.Temps["herms"] = resHerms.Temperature

	resMash, err := b.MashSensor.Get(ctx, &model.GetRequest{})
	if err != nil {
		return false, err
	}
	tps.Temps["mash"] = resMash.Temperature

	go func() {
		utils.LogIfError(&utils.DefualtPrinter{}, b.Logger.InsertTemperature(ctx, tps), "logging temp")
	}()

	scheme := b.Scheme.GetMash()
	// If the mash is greater than target, always turn off.
	if resMash.Temperature > scheme.MashMaxTemp {
		return false, nil
	}
	// If the boil is less than target, turn on.
	if resBoil.Temperature < scheme.BoilMinTemp {
		return true, nil
	}
	// If the herms is less than target, turn on
	if resHerms.Temperature < scheme.HermsMinTemp {
		return true, nil
	}
	// If the mash is less than target, turn on iff herms isn't too high.
	if resMash.Temperature < scheme.MashMinTemp {
		return resHerms.Temperature < scheme.HermsMaxTemp, nil
	}
	return false, nil

}

// StartRunLoop starts Brewery.Run in the background on a loop.
func (b *Brewery) StartRunLoop() {
	err := utils.RunLoop(5*time.Hour, 5*time.Second, b.run)
	if err != nil {
		utils.LogError(nil, err, "error running run")
	}
}

func (b *Brewery) run() error {
	b.mux.RLock()
	defer b.mux.RUnlock()
	if b.Scheme == nil {
		utils.Print("[Brewery.Run] No scheme present")
		return nil
	}
	config := b.Scheme
	switch sch := config.Scheme.(type) {
	case *model.ControlScheme_Mash_:
		utils.Print(fmt.Sprintf("Mashing: %+v", sch.Mash))
		return b.mash()
	case *model.ControlScheme_Boil_:
		return b.elementOn(context.Background()) // Toggle for one hour.
	}
	return nil
}

func (b *Brewery) elementOn(ctx context.Context) (err error) {
	_, err = b.Element.On(ctx, &model.OnRequest{})
	if err != nil {
		utils.LogError(nil, err, "encountered error turning coil on")
	}
	return err
}

func (b *Brewery) elementOff(ctx context.Context) error {
	var err error
	for i := 0; i < 3; i++ {
		_, err = b.Element.Off(ctx, &model.OffRequest{})
		if err == nil {
			return err
		}
	}
	return err
}

func (b *Brewery) mash() error {
	ctx := context.Background()
	on, err := b.mashThermOn(ctx)
	if err != nil {
		utils.Print(fmt.Sprintf("[Brewery.Run] MashThemOnErr: %s", err))
		return err
	}
	if !on {
		return b.elementOff(ctx)
	}
	return b.elementOn(ctx)
}
