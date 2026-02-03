package config

import (
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/goccy/go-yaml"
)

var (
	ErrNoConfigDir  = errors.New("configuration file directory cannot be determined")
	ErrBoxNotFound  = errors.New("box not found")
	ErrEmptyBoxName = errors.New("box name cannot be empty")
)

type Config struct {
	DefaultBox string         `json:"default_box,omitempty"`
	Boxes      map[string]Box `json:"boxes"`
}

func (c *Config) Box(boxName string) (Box, error) {
	boxName = cmp.Or(boxName, c.DefaultBox)

	if boxName == "" {
		return Box{}, ErrEmptyBoxName
	}

	box, ok := c.Boxes[boxName]
	if !ok {
		return Box{}, ErrBoxNotFound
	}

	return box, nil
}

type Box struct {
	People       []string `json:"people,omitempty"`
	Organization string   `json:"organization,omitempty"`
}

func Load() (Config, error) {
	var empty Config

	_, file, err := Location()
	if err != nil {
		return empty, err
	}

	data, err := os.ReadFile(file)
	if errors.Is(err, fs.ErrNotExist) {
		return empty, nil
	} else if err != nil {
		return empty, fmt.Errorf("reading file: %w", err)
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return empty, err
	}

	return config, nil
}

func Location() (dir, file string, err error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", "", ErrNoConfigDir
	}

	dir = path.Join(configDir, "gh-box")
	file = path.Join(dir, "config.yml")
	return dir, file, nil
}
