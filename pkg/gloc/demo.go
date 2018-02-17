package gloc

import "fmt"

func Demo(username, password string) error {
	client, err := NewLocationClient(username, password)
	if err != nil {
		return err
	}

	locations, err := client.Get()
	if err != nil {
		return err
	}

	for _, location := range locations {
		fmt.Println(location.String())
	}
	return nil
}
