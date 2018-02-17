package gloc

import (
	"testing"
)

var tta Location = Location{"KTTA", 35.5825, -79.1013333}
var w95 Location = Location{"W95", 35.1011667, -75.966}
var rdu Location = Location{"KRDU", 35.8776667, -78.7875}
var anc Location = Location{"PANC", 61.1741667, -149.9981667}
var lhd Location = Location{"PALH", 61.1816667, -149.9665}

type TCase struct {
	Distance float64
	Point1   Location
	Point2   Location
}

// correct values obtained from http://edwilliams.org/gccalc.htm
// using the spherical model
var testCases []TCase = []TCase{
	TCase{1.89, anc, lhd},
	TCase{43.3264, tta, rdu},
	TCase{289.1693, tta, w95},
	TCase{5585.1889, anc, rdu},
}

func TestDistance(t *testing.T) {
	for _, c := range testCases {
		low, high := halfPercent(c.Distance)
		d := c.Point1.Distance(c.Point2)

		// half-percent precision is good enough when dealing with cell phone
		// self-reported locations
		if low > d || high < d {
			t.Errorf("%f not within .5%% of %f", d, c.Distance)
		}

		// result should be the same both directions
		reverse := c.Point2.Distance(c.Point1)
		if d != reverse {
			t.Errorf("%f not equal to reverse calculation %f", d, reverse)
		}
	}
}

func halfPercent(x float64) (float64, float64) {
	return x * float64(.995), x * float64(1.005)
}
