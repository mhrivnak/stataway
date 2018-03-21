package google

import (
	"errors"
	"fmt"
	"github.com/mhrivnak/stataway/pkg/engine"
	"github.com/mhrivnak/stataway/pkg/gloc"
	"os"
	"strconv"
	"time"
)

type GDetector struct {
	home          bool
	client        *gloc.LocationClient
	triggerC      chan engine.Trigger
	homeLocation  gloc.Location
	innerDistance float64
	outerDistance float64
}

func New(triggerC chan engine.Trigger) (*GDetector, error) {
	username := os.Getenv("GOOGLE_USERNAME")
	password := os.Getenv("GOOGLE_PASSWORD")
	if username == "" || password == "" {
		return nil, errors.New("Google username and password must be provided.")
	}

	latS := os.Getenv("HOME_LATITUDE")
	lonS := os.Getenv("HOME_LONGITUDE")
	innerS := os.Getenv("GOOGLE_INNER_KM")
	outerS := os.Getenv("GOOGLE_OUTER_KM")
	if latS == "" || lonS == "" {
		return nil, errors.New("Home latitude and longitude must be provided.")
	}
	if innerS == "" || outerS == "" {
		return nil, errors.New("Inner and outer distances must be provided.")
	}
	lat, err := strconv.ParseFloat(latS, 64)
	if err != nil {
		return nil, fmt.Errorf("Could not parse %s into latitude", latS)
	}
	lon, err := strconv.ParseFloat(lonS, 64)
	if err != nil {
		return nil, fmt.Errorf("Could not parse %s into longitude", latS)
	}
	inner, err := strconv.ParseFloat(innerS, 64)
	if err != nil {
		return nil, fmt.Errorf("Could not parse %s into inner distance", innerS)
	}
	outer, err := strconv.ParseFloat(outerS, 64)
	if err != nil {
		return nil, fmt.Errorf("Could not parse %s into outer distance", outerS)
	}

	client, err := gloc.NewLocationClient(username, password)
	if err != nil {
		return nil, err
	}

	d := &GDetector{
		home:     true, // this will get set by Init() below
		client:   client,
		triggerC: triggerC,
		homeLocation: gloc.Location{
			Name:      "Home",
			Latitude:  lat,
			Longitude: lon,
		},
		innerDistance: inner,
		outerDistance: outer,
	}
	err = d.Init()
	return d, err
}

// Init uses current locations to determine if the initial state should be
// "home" or "away". This is important so that the detector does not trigger
// unneccessarily during startup.
func (d *GDetector) Init() error {
	minDist, err := d.smallestDistance()
	if err != nil {
		return err
	}
	// Assume if everyone is outside inner ring, that qualifies as "away" for
	// startup purposes.
	d.home = (minDist < d.innerDistance)
	return nil
}

// smallestDistance retrieves current locations and returns the smallest
// distance from one of those locations to the home location.
func (d *GDetector) smallestDistance() (float64, error) {
	locations, err := d.client.Get()
	if err != nil {
		return 0, err
	}

	var minDist float64 = d.outerDistance + 100 // an arbitrary value way above the threshold
	for _, loc := range locations {
		dist := loc.Distance(d.homeLocation)
		if dist < minDist {
			minDist = dist
		}
		fmt.Printf("%s is %f km from home.\n", loc.Name, dist)
	}
	return minDist, nil
}

// checkLocations sends a trigger for "home" or "away" if it determines, by
// using current location data, that a state change has occurred.
func (d *GDetector) checkLocations() {
	minDist, err := d.smallestDistance()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if minDist < d.innerDistance {
		if d.home == false {
			d.home = true
			d.triggerC <- engine.Trigger{true, "google", "inside home area"}
		}
	} else if minDist > d.outerDistance && d.home == true {
		d.home = false
		d.triggerC <- engine.Trigger{false, "google", "outside home area"}
	}
}

// Run initiates a loop where every 30 seconds the locations get checked.
func (d *GDetector) Run() {
	ticker := time.NewTicker(time.Second * 30)

	for _ = range ticker.C {
		d.checkLocations()
	}
}
