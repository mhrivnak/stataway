package venstar

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

// NewThermostat returns a Thermostat object ready to use. It depends on the
// environment variable VENSTAR_URL being set.
func NewThermostat() (*Thermostat, error) {
	host := os.Getenv("VENSTAR_URL")
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	client := http.Client{Timeout: time.Second * 30}

	return &Thermostat{u, &client}, nil
}

// Thermostat holds data about a physical thermostat on the network.
type Thermostat struct {
	url    *url.URL
	client *http.Client
}

// Set sets the thermostat's away mode
func (t *Thermostat) Set(home bool) error {
	path, err := url.Parse("settings")
	if err != nil {
		return err
	}
	u := t.url.ResolveReference(path)
	var away string
	if home {
		away = "0"
	} else {
		away = "1"
	}

	v := url.Values{}
	v.Set("away", away)

	resp, err := t.client.PostForm(u.String(), v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	result := Result{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}
	return result.OK()
}

// Home queries the thermostat and returns a bool indicating whether it is in
// "home" mode or not.
func (t *Thermostat) Home() (bool, error) {
	path, err := url.Parse("query/info")
	if err != nil {
		return false, err
	}
	resp, err := t.client.Get(t.url.ResolveReference(path).String())
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	i := Info{}
	err = json.Unmarshal(data, &i)
	if err != nil {
		return false, err
	}
	return i.Home()
}

// Info holds response data from a GET request to the API's "query/info"
// endpoint.
type Info struct {
	Away   int
	Error  bool
	Reason string
}

// Home returns a bool indicating whether the thermostat is in "home" mode or
// not.
func (i Info) Home() (bool, error) {
	if i.Error {
		message := fmt.Sprintf("thermostat returned error when queried: %s", i.Reason)
		return false, errors.New(message)
	}

	switch i.Away {
	case 0:
		return true, nil
	case 1:
		return false, nil
	default:
		message := fmt.Sprintf("thermostat returned invalid away state: %d", i.Away)
		return false, errors.New(message)
	}
}

// Result holds response data from a POST request to the API's "settings"
// endpoint.
type Result struct {
	Success bool
	Error   bool
	Reason  string
}

// OK returns nil if the response was ok, or an error if there was any problem
// with the response.
func (r Result) OK() error {
	if r.Error {
		message := fmt.Sprintf("thermostat returned error when being set: %s", r.Reason)
		return errors.New(message)
	}
	if r.Success == r.Error {
		return errors.New("invalid response from thermostat. success and failure have same value.")
	}
	return nil
}
