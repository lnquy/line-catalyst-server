package bot

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/lnquy/line-catalyst-server/pkg/aqi"
)

func (c *Catalyst) aqi(cmdArgs []string, replyTo string) error {
	city := c.conf.AQI.City
	if len(cmdArgs) >= 2 {
		city = cmdArgs[1]
	}
	w, err := aqi.GetAQIInfo(city, c.conf.AQI.Token)
	if err != nil {
		c.replyTo(replyTo, fmt.Sprintf("Sorry. Cannot get Air Quality Index for %s :(", strings.Title(city)))
		return errors.Wrapf(err, "failed to get Air Quality Index")
	}
	c.replyTo(replyTo, w)
	return nil
}
