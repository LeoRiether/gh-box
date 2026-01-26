package config

import (
	"fmt"
	"os"
)

type ConfigPathCmd struct{}

func (c *ConfigPathCmd) Run() error {
	dir, file, err := Location()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0777); err != nil {
		return fmt.Errorf("mkdir %q: %w", dir, err)
	}

	fmt.Println(file)
	return nil
}
