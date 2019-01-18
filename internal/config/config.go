package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type (
	Config struct {
		Port     string   `json:"port" required:"true"`
		LogLevel string   `json:"log_level" default:"info"`
		Bot      Bot      `json:"bot"`
		Database Database `json:"database"`
	}

	Bot struct {
		Secret      string      `json:"secret" required:"true"`
		Token       string      `json:"token" required:"true"`
		Translation Translation `json:"translation"`
		Weather     Weather     `json:"weather"`
		AQI         AQICN       `json:"aqi"`
		Joke        Joke        `json:"joke"`
	}

	Translation struct {
		SourceLang string `json:"source_language" default:"auto"`
		TargetLang string `json:"target_language" default:"en"`
	}

	Weather struct {
		Type        string         `json:"type" default:"openweathermap"`
		OpenWeather OpenWeatherMap `json:"open_weather"`
	}

	OpenWeatherMap struct {
		City  string `json:"city"`
		Token string `json:"token"`
	}

	AQICN struct {
		City  string `json:"city" default:"ho-chi-minh-city"`
		Token string `json:"token"`
	}

	Database struct {
		Type    string  `json:"type" required:"true" default:"mongodb"`
		MongoDB MongoDB `json:"mongo_db"`
	}

	Joke struct {
		Folder string `json:"folder" default:"_misc"`
	}

	MongoDB struct {
		URI string `json:"uri" envconfig:"URI"`
	}
)

func LoadEnvConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.Wrapf(err, "failed to read configurations from environment")
	}

	// Custom handling for MongoDB URI as Heroku store this value in MONGODB_URI
	if cfg.Database.Type == "mongodb" && cfg.Database.MongoDB.URI == "" {
		uri := os.Getenv("MONGODB_URI")
		if uri == "" {
			return nil, fmt.Errorf("unable to load MongoDB URI from environment variable")
		}
		cfg.Database.MongoDB.URI = uri
	}
	return &cfg, nil
}
