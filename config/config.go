package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Telegram Telegram `yaml:"telegram"`
	}

	Telegram struct {
		Token         string `yaml:"token"`
		UpdateTimeout int    `yaml:"updateTimeout"`
		Debug         bool   `yaml:"debug"`
	}
)

func LoadConfigFromFile(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to load config from \"%s\" - %s", file, err)
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file - %s", err)
	}

	return &config, nil
}
