package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Mongodb *Mongodb `yaml:"mongodb"`
	Redis   *Redis   `yaml:"redis"`
}

type Mongodb struct {
	URI string `yaml:"uri"`
}

type Redis struct {
	Addrs    []string `yaml:"addrs"`
	Password string   `yaml:"password"`
	DB       int      `yaml:"db"`
}

var (
	configFile = "./config/config.yaml"
)

func GetConfig() (*Config, error) {
	config := &Config{}

	// Read YAML file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML data into config struct
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
