package config

import (
	"fmt"

	"github.com/vrischmann/envconfig"
)

type (
	Config struct {
		Environment string `envconfig:"default=development"`
		// TODO: to support multiple channels, we might want to change this to map[string]string but not for now...
		Repos     map[string]bool `envconfig:"-"`
		ChannelID string          `envconfig:"-"`
	}
)

var (
	// NOTE: keep repos sorted alphabetically
	repos = map[string]bool{
		"indexer":   true,
		"mobileapi": true,
		"the-algo":  true,
	}
	configs = map[string]Config{
		"development": {
			Repos:     repos,
			ChannelID: "G78DB10F5",
		},
		"staging": {
			Repos:     repos,
			ChannelID: "G78DB10F5",
		},
		"production": {
			ChannelID: "C17P1LRD1",
		},
	}
)

func Load(env string) (Config, error) {
	config, ok := configs[env]
	if !ok {
		return config, fmt.Errorf("Unknown environment: %s", env)
	}

	err := envconfig.Init(&config)

	return config, err
}
