package main

import (
	"github.com/mhrivnak/stataway/pkg/detectors/gps"
	"github.com/mhrivnak/stataway/pkg/engine"
)

func main() {
	stateC := make(chan engine.State)
	triggerC := make(chan engine.Trigger)

	gpsD, err := gps.New(engine.State{true, false}, stateC, triggerC)
	if err != nil {
		panic(err.Error())
	}

	go gpsD.Run()

	err = engine.Run([]chan engine.State{stateC}, triggerC)
	if err != nil {
		panic(err.Error())
	}

	return
}
