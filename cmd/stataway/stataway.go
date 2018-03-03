package main

import (
	"github.com/mhrivnak/stataway/pkg/detectors/gps"
	"github.com/mhrivnak/stataway/pkg/engine"
)

func main() {
	triggerC := make(chan engine.Trigger)

	gpsD, err := gps.New(engine.State{true, false}, triggerC)
	if err != nil {
		panic(err.Error())
	}

	go gpsD.Run()

	err = engine.Run(triggerC)
	if err != nil {
		panic(err.Error())
	}

	return
}
