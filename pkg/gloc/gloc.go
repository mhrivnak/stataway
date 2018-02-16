package gloc

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"strings"
)

const locationURL = "https://www.google.com/maps/preview/locationsharing/read?gl=en&pb=%211m7%218m6%211m3%211i14%212i8413%213i5385%212i6%213x4095%212m3%211e0%212sm%213i407105169%213m7%212sen%215e1105%2112m4%211e68%212m2%211sset%212sRoadmap%214e1%215m4%211e4%218m2%211e0%211e1%216m9%211e12%212i2%2126m1%214b1%2130m1%211f1.3953487873077393%2139b1%2144e1%2150e0%2123i4111425&authuser=0&hl=en"

func Get(username, password string) error {
	client, err := login(username, password)
	if err != nil {
		return err
	}

	resp, err := client.Get(locationURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var hbuf bytes.Buffer
	_ = resp.Header.Write(&hbuf)
	fmt.Println(hbuf.String())

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	sdata := strings.SplitN(string(b), "'", 2)
	if len(sdata) != 2 {
		return errors.New("separator not found")
	}

	positions, err := jsonToLocation(sdata[1])
	if err != nil {
		return err
	}

	for _, position := range positions {
		fmt.Println(position.String())
	}
	return nil
}

func jsonToLocation(data string) ([]Location, error) {
	var ret []Location
	parsed := gjson.Parse(data)
	devices := parsed.Get("0").Array()
	for _, device := range devices {
		position := Location{
			Name:      device.Get("0.3").Str,
			Latitude:  device.Get("1.1.2").Num,
			Longitude: device.Get("1.1.1").Num,
		}
		if position.Latitude == 0 || position.Longitude == 0 {
			return ret, errors.New("could not find value in json")
		}
		ret = append(ret, position)
	}
	return ret, nil
}
