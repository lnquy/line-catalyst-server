package bot

import (
	"fmt"

	"github.com/line/line-bot-sdk-go/linebot"
)

const defaultMsgLimit = 5

func (c *Catalyst) translate(cmdArgs []string, replyTo string, isReplyToUser bool) error {
	_, err := c.bot.PushMessage(replyTo, linebot.NewTextMessage(fmt.Sprintf("Translate: %v", cmdArgs))).Do()
	return err
}
