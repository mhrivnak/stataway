package gloc

import "fmt"

type Location struct {
	Name      string
	Latitude  float64
	Longitude float64
}

func (l Location) String() string {
	return fmt.Sprintf("Lat: %f, Lon: %f, Name: %s", l.Latitude, l.Longitude, l.Name)
}

func (l Location) Distance(other Location) float64 {
	return 0
}
