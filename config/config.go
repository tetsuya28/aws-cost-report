package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	SlackToken   string `required:"true" envconfig:"SLACK_TOKEN"`
	SlackChannel string `required:"true" envconfig:"SLACK_CHANNEL"`
	Language     string `required:"true" envconfig:"LANGUAGE" default:"ja"`
}

func New() (Config, error) {
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		return Config{}, err
	}
	return config, nil
}
