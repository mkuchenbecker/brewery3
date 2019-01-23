package servers

import (
	"context"
	"fmt"
	"sync"
	"time"

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

	Element model.SwitchClient
}

// Control implements the Brewery.Control method. Given a ControlScheme
// The brewery will then attempt to follow the scheme.
func (b *Brewery) Control(ctx context.Context,
	req *model.ControlRequest) (res *model.ControlResponse, err error) {
	utils.Print("Recieved control request")
	b.replaceConfig(req.Scheme)
	return &model.ControlResponse{}, nil
}

func (b *Brewery) replaceConfig(scheme *model.ControlScheme) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Scheme = scheme
}

func (b *Brewery) getTempConstraints() ([]constraint, error) {
	resBoil, err := b.BoilSensor.Get(context.Background(), &model.GetRequest{})
	if err != nil {
		return []constraint{}, err
	}
	resHerms, err := b.HermsSensor.Get(context.Background(), &model.GetRequest{})
	if err != nil {
		return []constraint{}, err
	}
	resMash, err := b.MashSensor.Get(context.Background(), &model.GetRequest{})
	if err != nil {
		return []constraint{}, err
	}

	utils.Print(fmt.Sprintf("Mash: %f | Boil %f | Herms %f",
		resMash.Temperature, resBoil.Temperature, resHerms.Temperature))

	return []constraint{
		{
			min: b.Scheme.GetMash().BoilMinTemp,
			max: b.Scheme.GetMash().BoilMaxTemp,
			val: resBoil.Temperature,
		},
		{
			min: b.Scheme.GetMash().HermsMinTemp,
			max: b.Scheme.GetMash().HermsMaxTemp,
			val: resHerms.Temperature,
		},
		{
			min: b.Scheme.GetMash().MashMinTemp,
			max: b.Scheme.GetMash().MashMaxTemp,
			val: resMash.Temperature,
		},
	}, nil
}

func (b *Brewery) mashThermOn() (on bool, err error) {
	constraints, err := b.getTempConstraints()
	if err != nil {
		return false, err
	}
	val := checkTempConstraints(constraints)
	return val < 0, nil
}

type constraint struct {
	min float64
	max float64
	val float64
}

func (c *constraint) check() int {
	if c.val < c.min {
		return -1
	}
	if c.val >= c.max {
		return 1
	}
	return 0
}

// Returns -1 if some val is too low, 0 if all are met, and 1 if val is too high.
func checkTempConstraints(constriants []constraint) int {
	for _, constriant := range constriants {
		if val := constriant.check(); val != 0 {
			return val
		}
	}
	return 0
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
	case *model.ControlScheme_Boil_:
		return b.boil()
	case *model.ControlScheme_Mash_:
		return b.mash()
	case *model.ControlScheme_Power_:
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

func (b *Brewery) boil() error {
	err := b.elementOn()
	return err
}

func (b *Brewery) mash() error {
	on, err := b.mashThermOn()
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
