package bot

import (
	"github.com/pkg/errors"

	"github.com/lnquy/line-catalyst-server/pkg/weather"
)

func (c *Catalyst) weather(cmdArgs []string, replyTo string) error {
	city := c.conf.Weather.OpenWeather.City
	if len(cmdArgs) >= 2 {
		city = cmdArgs[1]
	}
	w, err := weather.GetWeatherInfo(city, c.conf.Weather.OpenWeather.Token)
	if err != nil {
		return errors.Wrapf(err, "failed to get weather info")
	}
	c.replyTo(replyTo, w)
	return nil
}
