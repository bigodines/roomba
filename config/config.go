package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	yaml "github.com/ghodss/yaml"
	"github.com/vrischmann/envconfig"
)

type (
	Config struct {
		Environment string `envconfig:"default=development"`
		Webhook     string `json:"web_hook" envconfig:"-"`
		// TODO: to support multiple channels, we might want to change this to map[string]string but not for now...
		Repos        map[string]bool   `json:"repos" envconfig:"-"`
		Countdown    map[string]string `json:"countdown" envconfig:"-"`
		ChannelID    string            `json:"channel_id" envconfig:"-"`
		Organization string            `json:"organization" envconfig:"-"`
	}
)

func Load(env string) (Config, error) {
	config, err := configFromFile(env)
	if err != nil {
		return config, err
	}
	config.Environment = env

	err = envconfig.Init(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// configurationFromFile reads configuration file for environment and return a Config struct
func configFromFile(env string) (Config, error) {
	env = strings.ToLower(env)
	var fname string
	var conf Config
	if _, err := os.Stat(fmt.Sprintf("config/%s.yml", env)); err == nil {
		fname = fmt.Sprintf("config/%s.yml", env)
	} else {
		fname = fmt.Sprintf("config/default.yml")
	}

	ymlFile, err := ioutil.ReadFile(fname)
	if err != nil {
		return conf, err
	}

	err = yaml.Unmarshal(ymlFile, &conf)
	if err != nil {
		return conf, err
	}

	return conf, nil

}
