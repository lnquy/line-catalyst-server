package bot

import (
	"github.com/pkg/errors"

	"github.com/lnquy/line-catalyst-server/pkg/weather"
)

func (c *Catalyst) weather(cmdArgs []string, replyTo string) error {
	w, err := weather.AQIInfo(c.conf.Weather.Aqicn.City, c.conf.Weather.Aqicn.Token)
	if err != nil {
		return errors.Wrapf(err, "failed to get AQICN weather info")
	}
	c.replyTo(replyTo, w)
	return nil
}
