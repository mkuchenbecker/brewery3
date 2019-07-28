package main

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	gomodel "github.com/mkuchenbecker/brewery3/brewery/model/gomodel"
)

func main() {

	settings := gomodel.BrewerySettings{
		Config: &gomodel.ServerConfig{
			Name: "brewery0",
			Addr: "localhost",
			Port: 9000,
		},
		Thermometers: []*gomodel.ServerConfig{
			&gomodel.ServerConfig{
				Name: "mash",
				Addr: "localhost",
				Port: 9110,
			},
			&gomodel.ServerConfig{
				Name: "boil",
				Addr: "localhost",
				Port: 9112,
			},
			&gomodel.ServerConfig{
				Name: "herms",
				Addr: "localhost",
				Port: 9111,
			},
		},
		Heater: &gomodel.ServerConfig{
			Name: "main",
			Addr: "localhost",
			Port: 9100,
		},
	}
	out, err := yaml.Marshal(settings)
	if err != nil {
		panic(err)
	}

	// file, err := os.Create("config.yaml")
	// if err != nil {
	// 	panic(err)
	// }
	// defer func() {
	// 	defErr := file.Close()
	// 	if err == nil {
	// 		err = defErr
	// 	}
	// }()
	err = ioutil.WriteFile("config.yaml", out, 775)
	if err != nil {
		panic(err)
	}
}
