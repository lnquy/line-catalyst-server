package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type (
	Config struct {
		Port     string `required:"true"`
		Bot      Bot
		Database Database
	}

	Bot struct {
		Secret string `required:"true"`
		Token  string `required:"true"`
	}

	Database struct {
		Type    string `required:"true" default:"mongodb"`
		MongoDB MongoDB
	}

	MongoDB struct {
		URI string `envconfig:"URI" required:"true"`
	}
)

func LoadEnvConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.Wrapf(err, "failed to read configurations from environment")
	}
	return &cfg, nil
}
