package engine

import (
	"fmt"
	"github.com/mhrivnak/stataway/pkg/thermostat/venstar"
)

func Run(triggerC chan Trigger) error {
	stat, err := venstar.NewThermostat()
	if err != nil {
		return err
	}

	for {
		t := <-triggerC
		fmt.Println(t.String())
		home, err := stat.Home()
		if err != nil {
			fmt.Println("Error getting thermostat state:", err.Error())
			continue
		}
		fmt.Println("Currently set to home:", home)
		if home == t.Home {
			fmt.Println("Thermostat home state already agrees with trigger")
		} else {
			err = stat.Set(t.Home)
			if err != nil {
				fmt.Println("Error setting thermostat state:", err.Error())
				continue
			}
			fmt.Println("Set thermostat home state to:", t.Home)
		}
	}
}

