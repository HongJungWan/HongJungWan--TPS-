package config

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type Config struct {
	DB struct {
		Database string
		URL      string
	}

	Kafka struct {
		ClientID string
		URL      string
	}
}

func NewConfig(path string) *Config {
	config := new(Config)

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = toml.NewDecoder(file).Decode(config)
	if err != nil {
		panic(err)
	}

	return config
}
