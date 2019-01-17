package bot

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"
)

func (c *Catalyst) weather(cmdArgs []string, replyTo string) error {
	_, err := c.bot.PushMessage(replyTo, linebot.NewTextMessage(fmt.Sprintf("Weather: %v", cmdArgs))).Do()
	return err
}
