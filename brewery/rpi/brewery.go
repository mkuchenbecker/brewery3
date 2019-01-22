package rpi

import (
	"context"
	"fmt"
	"sync"
	"time"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
	"github.com/mkuchenbecker/brewery3/brewery/utils"
)

type Brewery struct {
	Scheme *model.ControlScheme
	mux    sync.RWMutex

	MashSensor  model.ThermometerClient
	BoilSensor  model.ThermometerClient
	HermsSensor model.ThermometerClient

	mashTemp  float64
	boilTemp  float64
	hermsTemp float64
	tempMux   sync.RWMutex

	tempsRead bool

	Element model.SwitchClient
}

func (c *Brewery) Control(ctx context.Context,
	req *model.ControlRequest) (res *model.ControlResponse, err error) {
	utils.Print("Recieved control request")
	c.ReplaceConfig(req.Scheme)
	return &model.ControlResponse{}, nil
}

func (c *Brewery) ReplaceConfig(scheme *model.ControlScheme) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.Scheme = scheme
}

func (c *Brewery) updateTemperatures() error {
	resBoil, err := c.BoilSensor.Get(context.Background(), &model.GetRequest{})
	if err != nil {
		return err
	}
	resHerms, err := c.HermsSensor.Get(context.Background(), &model.GetRequest{})
	if err != nil {
		return err
	}
	resMash, err := c.MashSensor.Get(context.Background(), &model.GetRequest{})
	if err != nil {
		return err
	}

	utils.Print(fmt.Sprintf("Mash: %f | Boil %f | Herms %f",
		resMash.Temperature, resBoil.Temperature, resHerms.Temperature))

	c.tempMux.Lock()
	defer c.tempMux.Unlock()
	c.boilTemp = resBoil.Temperature
	c.mashTemp = resMash.Temperature
	c.hermsTemp = resHerms.Temperature
	c.tempsRead = true
	return nil
}

func (c *Brewery) getTempConstraints() ([]Constraint, error) {

	c.tempMux.RLock()
	updateLive := !c.tempsRead
	c.tempMux.RUnlock()

	if updateLive {
		err := c.updateTemperatures()
		if err != nil {
			return []Constraint{}, err
		}
	} else {
		go c.updateTemperatures()
	}

	c.tempMux.RLock()
	defer c.tempMux.RUnlock()

	return []Constraint{
		{
			min: c.Scheme.GetMash().BoilMinTemp,
			max: c.Scheme.GetMash().BoilMaxTemp,
			val: c.boilTemp,
		},
		{
			min: c.Scheme.GetMash().HermsMinTemp,
			max: c.Scheme.GetMash().HermsMaxTemp,
			val: c.hermsTemp,
		},
		{
			min: c.Scheme.GetMash().MashMinTemp,
			max: c.Scheme.GetMash().MashMaxTemp,
			val: c.mashTemp,
		},
	}, nil
}

func (c *Brewery) mashThermOn() (on bool, err error) {
	constraints, err := c.getTempConstraints()
	if err != nil {
		return false, err
	}
	val := checkTempConstraints(constraints)
	return val < 0, nil
}

type Constraint struct {
	min float64
	max float64
	val float64
}

func (c *Constraint) Check() int {
	if c.val < c.min {
		return -1
	}
	if c.val >= c.max {
		return 1
	}
	return 0
}

// Returns -1 if some val is too low, 0 if all are met, and 1 if val is too high.
func checkTempConstraints(constriants []Constraint) int {
	for _, constriant := range constriants {
		if val := constriant.Check(); val != 0 {
			return val
		}
	}
	return 0
}

func (c *Brewery) ElementOff() error {
	var err error
	for i := 0; i < 3; i++ {
		_, err = c.Element.Off(context.Background(), &model.OffRequest{})
		if err == nil {
			return err
		}
	}
	return err
}

func (b *Brewery) RunLoop() error {
	ttl := 5
	for {
		utils.Print("[RunLoop]")
		err := b.Run()
		if err != nil {
			utils.Print(fmt.Sprintf("[RunLoop] %s", err.Error()))
			time.Sleep(time.Duration(ttl) * time.Second)
		}
		err = nil
	}
}

func (b *Brewery) Run() error {
	b.mux.RLock()
	defer b.mux.RUnlock()
	if b.Scheme == nil {
		utils.Print("[Brewery.Run] No scheme present")
		return fmt.Errorf("no scheme present")
	}
	config := b.Scheme
	switch sch := config.Scheme.(type) {
	case *model.ControlScheme_Boil_:
		err := b.ElementOn()
		return err
	case *model.ControlScheme_Mash_:
		on, err := b.mashThermOn()
		if err != nil {
			utils.Print(fmt.Sprintf("[Brewery.Run] MashThemOnErr: %s", err))
			return err
		}
		if !on {
			return b.ElementOff()
		}
		return b.ElementOn()
	case *model.ControlScheme_Power_:
		return b.ElementPowerLevel(int(sch.Power.PowerLevel)) // Toggle for one hour.
	default:
	}
	return nil
}

func (b *Brewery) ElementOn() (err error) {
	_, err = b.Element.On(context.Background(), &model.OnRequest{})
	if err != nil {
		print(fmt.Sprintf("encountered error turning coil on: %+v", err))
	}
	return err
}

func (b *Brewery) ElementPowerLevel(powerLevel int) error {
	if powerLevel < 1 {
		return b.ElementOff()
	}
	if powerLevel > 100 {
		return b.ElementOff()
	}
	if powerLevel == 100 {
		return b.ElementOn()
	}
	interval := 2
	delay := time.Duration(powerLevel / 100 * interval)
	return b.elementPowerLevelToggle(delay, 10*time.Second, time.Duration(interval)*time.Second)
}

func (b *Brewery) elementPowerLevelToggle(delay time.Duration, ttl time.Duration, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	quit := make(chan bool)
	resErr := make(chan error)

	go func() {
		for {
			select {
			case <-ticker.C:
				timer := time.NewTimer(delay)
				utils.Print(".")
				err := b.ElementOn()
				if err != nil {
					resErr <- err
					return
				}

				<-timer.C
				err = b.ElementOff()
				if err != nil {
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
