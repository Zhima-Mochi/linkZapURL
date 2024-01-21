package config

import (
	"os"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Mongodb *Mongodb `yaml:"mongodb"`
	Redis   *Redis   `yaml:"redis"`
}

type Mongodb struct {
	URI         string `yaml:"uri"`
	Database    string `yaml:"database"`
	MaxPoolSize uint64 `yaml:"max_pool_size"`
	USERNAME    string `yaml:"username"`
	PASSWORD    string `yaml:"password"`
}

type Redis struct {
	Addrs    []string `yaml:"addrs"`
	Password string   `yaml:"password"`
	DB       int      `yaml:"db"`
}

var (
	configFile = "config/config.yaml"
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

func (m *Mongodb) GetClientOptions() []*options.ClientOptions {
	opts := []*options.ClientOptions{}

	if m.URI != "" {
		opts = append(opts, options.Client().ApplyURI(m.URI))
	}

	if m.MaxPoolSize != 0 {
		opts = append(opts, options.Client().SetMaxPoolSize(m.MaxPoolSize))
	}

	// Low Write, High Read
	opts = append(opts, options.Client().SetReadConcern(readconcern.Local()))
	opts = append(opts, options.Client().SetWriteConcern(writeconcern.Majority()))
	opts = append(opts, options.Client().SetReadPreference(readpref.Nearest()))

	return opts
}
