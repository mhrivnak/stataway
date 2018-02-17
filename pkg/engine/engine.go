package engine

import "fmt"

func Run(stateChannels []chan State, triggerC chan Trigger) error {
	for {
		t := <-triggerC
		fmt.Println(t.String())
	}
}
