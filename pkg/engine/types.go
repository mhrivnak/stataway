package engine

import "fmt"

type Trigger struct {
	State        string
	DetectorName string
	Reason       string
}

func (t Trigger) String() string {
	return fmt.Sprintf("Trigger{State: %s, DetectorName: %s, Reason: %s}", t.State, t.DetectorName, t.Reason)
}

type State struct {
	Home   bool
	Paused bool
}

type Detector interface {
	Run(State, chan State, chan Trigger) error
}
