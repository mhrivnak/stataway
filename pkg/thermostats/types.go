package thermostats

type Thermostat interface {
	Set(bool) error
	Home() (bool, error)
}
