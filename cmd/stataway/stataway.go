package main

import (
	"github.com/mhrivnak/stataway/pkg/detectors/google"
	"github.com/mhrivnak/stataway/pkg/engine"
	"github.com/mhrivnak/stataway/pkg/thermostats/venstar"
)

func main() {
	triggerC := make(chan engine.Trigger)

	stat, err := venstar.NewThermostat()
	if err != nil {
		panic(err.Error())
	}

	googleD, err := google.New(triggerC)
	if err != nil {
		panic(err.Error())
	}

	go googleD.Run()

	// this will block unless something goes wrong
	err = engine.Run(stat, triggerC)
	if err != nil {
		panic(err.Error())
	}

	return
}
