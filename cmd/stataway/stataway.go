package main

import (
	"github.com/mhrivnak/stataway/pkg/detectors/gps"
	"github.com/mhrivnak/stataway/pkg/engine"
	"github.com/mhrivnak/stataway/pkg/thermostat/venstar"
)

func main() {
	triggerC := make(chan engine.Trigger)

	stat, err := venstar.NewThermostat()
	if err != nil {
		panic(err.Error())
	}

	gpsD, err := gps.New(engine.State{true, false}, triggerC)
	if err != nil {
		panic(err.Error())
	}

	go gpsD.Run()

	err = engine.Run(stat, triggerC)
	if err != nil {
		panic(err.Error())
	}

	return
}
