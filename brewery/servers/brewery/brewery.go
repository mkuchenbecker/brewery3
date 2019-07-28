package brewery

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/mkuchenbecker/brewery3/brewery/logger"
	"github.com/mkuchenbecker/brewery3/brewery/utils"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
)

type SensorName string

// Brewery is the central controller for various heating elements, sensors, and switches.
// It coordinates all actions of the brewery.
type Brewery struct {
	Scheme *model.ControlScheme
	mux    sync.RWMutex

	Temperatures map[SensorName]model.ThermometerClient
	Log          logger.Logger

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

func (b *Brewery) getTemperaturePointSink(ctx context.Context) (*utils.TemperaturePointSink, error) {

	tps := utils.NewTemperaturePointSink()

	for sensorName, sensor := range b.Temperatures {
		response, err := sensor.Get(ctx, &model.GetRequest{})
		if err != nil {
			return nil, err
		}
		tps.Temps[string(sensorName)] = response.Temperature
	}

	b.Log.InsertTemperature(ctx, tps.ToInterface())

	return tps, nil
}

func (b *Brewery) mashThermOn(ctx context.Context, ctrl *model.ControlScheme_Mash) (on bool, err error) {
	tps, err := b.getTemperaturePointSink(ctx)

	for _, constraint := range ctrl.Constraints {
		val, err := tps.Check(constraint.Key, constraint.Min, constraint.Max)
		if err != nil {
			return false, err
		}
		if val < 0 {
			return true, nil
		}
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
		return b.mash(context.Background(), sch.Mash)
	case *model.ControlScheme_Power_:
		utils.Print(fmt.Sprintf("Power level: %f", sch.Power.PowerLevel))
		return b.powerLevel(int(sch.Power.PowerLevel)) // Toggle for one hour.
	}
	return nil
}

func (b *Brewery) elementOn() (err error) {
	_, err = b.Element.On(context.Background(), &model.OnRequest{})
	if err != nil {
		utils.LogError(nil, err, "encountered error turning coil on")
	}
	return err
}

func (b *Brewery) elementToggle(intervalMs int64) (err error) {
	_, err = b.Element.ToggleOn(context.Background(), &model.ToggleOnRequest{IntervalMs: intervalMs})
	if err != nil {
		utils.LogError(nil, err, "encountered error toggling coil")
	}
	return err
}

func (b *Brewery) elementOff() error {
	var err error
	for i := 0; i < 3; i++ {
		_, err = b.Element.Off(context.Background(), &model.OffRequest{})
		if err == nil {
			return err
		}
	}
	return err
}

func (b *Brewery) mash(ctx context.Context, ctrl *model.ControlScheme_Mash) error {
	on, err := b.mashThermOn(ctx, ctrl)
	if err != nil {
		utils.Print(fmt.Sprintf("[Brewery.Run] MashThemOnErr: %s", err))
		return err
	}
	if !on {
		return b.elementOff()
	}
	return b.elementOn()
}

func (b *Brewery) powerLevel(powerLevel int) error {
	if powerLevel < 1 {
		return b.elementOff()
	}
	if powerLevel > 100 {
		return b.elementOff()
	}
	if powerLevel == 100 {
		return b.elementOn()
	}
	interval := int64(2)
	intervalMs := int64((1000*powerLevel)/100) * interval
	return b.powerToggle(intervalMs, time.Duration(interval)*time.Second, 10*time.Second)
}

func (b *Brewery) powerToggle(intervalMs int64, loopInterval time.Duration, ttl time.Duration) (err error) {
	ticker := time.NewTicker(loopInterval)
	quit := make(chan bool)
	resErr := make(chan error)

	defer utils.DeferErrReturn(b.elementOff, &err)
	go func() {
		for {
			select {
			case <-ticker.C:
				utils.Print(".")
				err := b.elementToggle(intervalMs)
				if err != nil {
					utils.LogError(nil, err, "element toggle error")
					resErr <- err
					return
				}
			case <-quit:
				ticker.Stop()
				resErr <- nil
				return
			}
		}
	}()

	go func() { // Make sure the process always exits.
		timer := time.NewTimer(ttl)
		<-timer.C
		quit <- true
	}()

	return <-resErr
}
