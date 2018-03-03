package engine

import "fmt"

type Trigger struct {
	Home         bool
	DetectorName string
	Reason       string
}

func (t Trigger) String() string {
	return fmt.Sprintf("Trigger{Home: %t, DetectorName: %s, Reason: %s}", t.Home, t.DetectorName, t.Reason)
}
