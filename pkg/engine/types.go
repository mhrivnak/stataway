package engine

import "fmt"

type Trigger struct {
	Home         bool
	DetectorName string
	Reason       string
}

func (t Trigger) String() string {
	return fmt.Sprintf("Trigger{Home: %s, DetectorName: %s, Reason: %s}", t.Home, t.DetectorName, t.Reason)
}

type State struct {
	Home   bool
	Paused bool
}

type Detector interface {
	Run(State, chan Trigger) error
}
