package utils

import (
	"fmt"
	"log"
)

type TemperaturePointSink struct {
	Temps map[string]float64
}

func NewTemperaturePointSink() *TemperaturePointSink {
	return &TemperaturePointSink{
		Temps: make(map[string]float64),
	}
}

func (tps *TemperaturePointSink) ToInterface() map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range tps.Temps {
		res[k] = interface{}(v)
	}
	return res
}

func (tps *TemperaturePointSink) Log() {
	for k, temp := range tps.Temps {
		Printf("Temp [%s]%f", k, temp)
	}
}

// 0 = within range.
// -1 = Below range
//  1 = above range
func (tps *TemperaturePointSink) Check(key string, min float64, max float64) (int, error) {
	val, ok := tps.Temps[key]
	if !ok {
		return 0, fmt.Errorf("key not found: %s", key)
	}
	if val < min {
		return -1, nil
	}
	if val >= max {
		return 1, nil
	}
	return 0, nil
}

func (tps *TemperaturePointSink) CheckMust(key string, min float64, max float64) int {
	res, err := tps.Check(key, min, max)
	if err != nil {
		log.Fatal(err)
	}
	return res
}
