package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	SlackToken   string `required:"true" envconfig:"SLACK_TOKEN"`
	SlackChannel string `required:"true" envconfig:"SLACK_CHANNEL"`
}

func New() (*Config, error) {
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		return nil, err
	}
	return &Config{}, nil
}
