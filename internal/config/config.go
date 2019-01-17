package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type (
	Config struct {
		Port     string   `json:"port" required:"true"`
		Bot      Bot      `json:"bot"`
		Database Database `json:"database"`
	}

	Bot struct {
		Secret string `json:"secret" required:"true"`
		Token  string `json:"token" required:"true"`
	}

	Database struct {
		Type    string  `json:"type" required:"true" default:"mongodb"`
		MongoDB MongoDB `json:"mongo_db"`
	}

	MongoDB struct {
		URI string `json:"uri" envconfig:"MONGODB_URI" required:"true"`
	}
)

func LoadEnvConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.Wrapf(err, "failed to read configurations from environment")
	}
	return &cfg, nil
}
