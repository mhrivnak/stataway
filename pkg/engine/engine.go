package engine

import (
	"fmt"
	"github.com/mhrivnak/stataway/pkg/thermostats"
)

func Run(stat thermostats.Thermostat, triggerC chan Trigger) error {

	for t := range triggerC {
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
	return nil
}
