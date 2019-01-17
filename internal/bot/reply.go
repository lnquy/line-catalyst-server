package bot

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/sirupsen/logrus"
)

func (c *Catalyst) replyTo(to string, message string) {
	_, err := c.bot.PushMessage(to, linebot.NewTextMessage(message)).Do()
	if err != nil {
		// Log only, do nothing if failed to reply for now
		logrus.Errorf("bot: failed to reply message to %s user: %v", to, err)
	}
}
