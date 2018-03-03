package gps

import (
	"errors"
	"fmt"
	"github.com/mhrivnak/stataway/pkg/engine"
	"github.com/mhrivnak/stataway/pkg/gloc"
	"os"
	"strconv"
	"time"
)

type GPSDetector struct {
	state         string
	home          bool
	client        *gloc.LocationClient
	triggerC      chan engine.Trigger
	homeLocation  gloc.Location
	innerDistance float64
	outerDistance float64
}

func New(triggerC chan engine.Trigger) (*GPSDetector, error) {
	username := os.Getenv("GOOGLE_USERNAME")
	password := os.Getenv("GOOGLE_PASSWORD")
	if username == "" || password == "" {
		return nil, errors.New("Google username and password must be provided")
	}

	latS := os.Getenv("HOME_LATITUDE")
	lonS := os.Getenv("HOME_LONGITUDE")
	if latS == "" || lonS == "" {
		return nil, errors.New("Home latitude and longitude must be provided")
	}
	lat, err := strconv.ParseFloat(latS, 64)
	if err != nil {
		return nil, fmt.Errorf("Could not parse %s into latitude", latS)
	}
	lon, err := strconv.ParseFloat(lonS, 64)
	if err != nil {
		return nil, fmt.Errorf("Could not parse %s into longitude", latS)
	}

	client, err := gloc.NewLocationClient(username, password)
	if err != nil {
		return nil, err
	}

	d := &GPSDetector{
		home:     true, // this will get set by Init() below
		client:   client,
		triggerC: triggerC,
		homeLocation: gloc.Location{
			Name:      "Home",
			Latitude:  lat,
			Longitude: lon,
		},
		innerDistance: 0.5,
		outerDistance: 0.7,
	}
	err = d.Init()
	return d, err
}

func (d *GPSDetector) Init() error {
	minDist, err := d.smallestDistance()
	if err != nil {
		return err
	}
	// Assume if everyone is outside inner ring, that qualifies as "away" for
	// startup purposes.
	d.home = (minDist < d.innerDistance)
	return nil
}

func (d *GPSDetector) smallestDistance() (float64, error) {
	locations, err := d.client.Get()
	if err != nil {
		return 0, err
	}

	var minDist float64 = 100 // an arbitrary value way above the threshold
	for _, loc := range locations {
		dist := loc.Distance(d.homeLocation)
		if dist < minDist {
			minDist = dist
		}
		fmt.Printf("%s is %f km from home.\n", loc.Name, dist)
	}
	return minDist, nil
}

func (d *GPSDetector) checkLocations() {
	minDist, err := d.smallestDistance()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if minDist < d.innerDistance {
		if d.home == false {
			d.home = true
			d.triggerC <- engine.Trigger{true, "gps", "inside home area"}
		}
	} else if minDist > d.outerDistance && d.home == true {
		d.home = false
		d.triggerC <- engine.Trigger{false, "gps", "outside home area"}
	}
}

func (d *GPSDetector) Run() {
	tickerC := time.NewTicker(time.Second * 30).C

	for {
		select {
		case <-tickerC:
			d.checkLocations()
		}
	}
}
