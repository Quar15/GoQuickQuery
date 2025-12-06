package config

import (
	"os"

	"github.com/quar15/qq-go/internal/database"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Connections []database.ConnectionData `yaml:"connections"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
