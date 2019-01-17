package bot

import "github.com/line/line-bot-sdk-go/linebot"

const helpMsg = `Usage: @metion [command] [args]

- Mention bot by: @Catalyst, @catalyst, @tr, @th, @en, :tr, :th or :en
- Commands: 
   + translate: Translate last N messages to English.
   + weather:   Report weather info.
   + help/?:    Show this help message.
`

func (c *Catalyst) help(replyTo string) error {
	_, err := c.bot.PushMessage(replyTo, linebot.NewTextMessage(helpMsg)).Do()
	return err
}
