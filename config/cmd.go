package config

import "fmt"

type ConfigPathCmd struct{}

func (c *ConfigPathCmd) Run() error {
	file, err := Location()
	if err != nil {
		return err
	}

	fmt.Println(file)
	return nil
}
