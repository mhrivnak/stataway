package gloc

import (
	"fmt"
	"math"
)

type Location struct {
	Name      string
	Latitude  float64
	Longitude float64
}

const p float64 = math.Pi / 180

// radius of Earth in km
const radius int = 6371

func (l Location) String() string {
	return fmt.Sprintf("Lat: %f, Lon: %f, Name: %s", l.Latitude, l.Longitude, l.Name)
}

// Distance - returns the great circle distance between two points in
// kilometers using the Haversine formula.
func (l Location) Distance(other Location) float64 {
	a := 0.5 - math.Cos((l.Latitude-other.Latitude)*p)/2 + math.Cos(other.Latitude*p)*math.Cos(l.Latitude*p)*(1-math.Cos((l.Longitude-other.Longitude)*p))/2
	return float64(2*radius) * math.Asin(math.Sqrt(a))
}
