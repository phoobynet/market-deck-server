package server

import (
	"github.com/pelletier/go-toml"
	"os"
)

type Config struct {
	ServerPort int `toml:"server_port"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var config Config

	err = toml.Unmarshal(file, &config)

	if err != nil {
		return nil, err
	}

	return &config, nil
}
