package rpi

import (
	"context"
	"fmt"
	"sync"
	"time"

	model "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
)

func print(s string) {
	fmt.Printf("%s %s\n", time.Now().Format(time.StampMilli), s)
}

type Brewery struct {
	scheme *model.ControlScheme
	mux    sync.RWMutex

	mashSensor  model.ThermometerClient
	boilSensor  model.ThermometerClient
	hermsSensor model.ThermometerClient

	element model.SwitchClient
}

func (c *Brewery) ReplaceConfig(scheme *model.ControlScheme) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.scheme = scheme
}

func (c *Brewery) getTempConstraints() ([]Constraint, error) {
	resBoil, err := c.boilSensor.Get(context.Background(), &model.GetRequest{})
	if err != nil {
		return []Constraint{}, err
	}
	resHerms, err := c.hermsSensor.Get(context.Background(), &model.GetRequest{})
	if err != nil {
		return []Constraint{}, err
	}
	resMash, err := c.mashSensor.Get(context.Background(), &model.GetRequest{})
	if err != nil {
		return []Constraint{}, err
	}

	return []Constraint{
		{
			min: c.scheme.GetMash().BoilMinTemp,
			max: c.scheme.GetMash().BoilMaxTemp,
			val: resBoil.Temperature,
		},
		{
			min: c.scheme.GetMash().HermsMinTemp,
			max: c.scheme.GetMash().HermsMaxTemp,
			val: resHerms.Temperature,
		},
		{
			min: c.scheme.GetMash().MashMinTemp,
			max: c.scheme.GetMash().MashMaxTemp,
			val: resMash.Temperature,
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
		_, err = c.element.Off(context.Background(), &model.OffRequest{})
		if err == nil {
			return err
		}
	}
	return err
}

func (b *Brewery) Run(ttlSec int) error {
	b.mux.RLock()
	defer b.mux.RUnlock()
	ttl := time.Duration(ttlSec) * time.Second
	config := b.scheme
	switch sch := config.Scheme.(type) {
	case *model.ControlScheme_Boil_:
		err := b.ElementOn(ttl)
		return err
	case *model.ControlScheme_Mash_:
		on, err := b.mashThermOn()
		if err != nil {
			return err
		}
		if !on {
			err := b.ElementOff()
			if err != nil {
				return err
			}
			return nil
		}
		b.ElementOn(ttl)
		return err
	case *model.ControlScheme_Power_:
		b.ElementPowerLevel(int(sch.Power.PowerLevel), ttlSec) // Toggle for one hour.
	default:
	}
	return nil
}

func (b *Brewery) ElementOn(ttl time.Duration) (err error) {
	defer func() {
		offErr := b.ElementOff()
		if offErr != nil || err != nil {
			err = fmt.Errorf("errors occured: '%s', '%s", offErr, err)
		}
	}()

	_, err = b.element.On(context.Background(), &model.OnRequest{})
	if err != nil {
		print(fmt.Sprintf("encountered error turning coil on: %+v", err))
		return err
	}
	timer := time.NewTimer(ttl)
	<-timer.C
	return err
}

func (b *Brewery) ElementPowerLevel(powerLevel int, ttlSeconds int) error {
	ttl := time.Duration(ttlSeconds) * time.Second
	if powerLevel < 1 {
		err := b.ElementOff()
		if err != nil {
			return err
		}
	}
	if powerLevel > 100 {
		err := b.ElementOff()
		if err != nil {
			return err
		}
	}
	if powerLevel == 100 {
		err := b.ElementOn(ttl)
		if err != nil {
			return err
		}
	}
	interval := 2
	delay := time.Duration(powerLevel / 100 * interval)
	return b.elementPowerLevelToggle(delay, ttl, time.Duration(interval)*time.Second)
}

func (b *Brewery) elementPowerLevelToggle(delay time.Duration, ttl time.Duration, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	quit := make(chan bool)
	resErr := make(chan error)

	go func() {
		for {
			select {
			case <-ticker.C:
				err := b.ElementOn(delay)
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
