package config

import (
	"log/slog"
	"os"
	"sync"

	"github.com/quar15/qq-go/internal/database"
	"gopkg.in/yaml.v3"
)

func init() {
	_, _ = Load()
}

type Config struct {
	Connections []database.ConnectionData `yaml:"connections"`
	Colors      colors                    `yaml:"colors,omitempty"`
}

var (
	cfg  *Config = &Config{}
	once sync.Once
	err  error
)

func Load() (*Config, error) {
	const connectionsConfigPath = "./config/gqq.yaml"
	const colorsConfigPath = "./config/colors.yaml"
	slog.Debug(
		"Trying to initialize config",
		slog.String("connectionsConfigPath", connectionsConfigPath),
		slog.String("colorsConfigPath", colorsConfigPath),
	)
	once.Do(func() {
		var data []byte
		data, err = os.ReadFile(connectionsConfigPath)
		if err != nil {
			return
		}

		conns := []database.ConnectionData{}
		if err = yaml.Unmarshal(data, &conns); err != nil {
			return
		}
		slog.Debug("Initialized connections from config", slog.Any("conns", conns))

		data, err = os.ReadFile(colorsConfigPath)
		if err != nil {
			return
		}

		colorsCfg := ColorsConfig{}
		if err = yaml.Unmarshal(data, &colorsCfg); err != nil {
			return
		}
		colorsCfg.MergeDefaults()
		slog.Debug("Initialized colors from config", slog.Any("colorsCfg", colorsCfg))

		cfg = &Config{
			Connections: conns,
			Colors:      colors{&colorsCfg},
		}
	})

	slog.Debug("Initialized config", slog.Any("config", cfg))

	return cfg, err
}

func Get() *Config {
	if cfg == nil {
		panic("Config not loaded: call Load first")
	}
	return cfg
}
