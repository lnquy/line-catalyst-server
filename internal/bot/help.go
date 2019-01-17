package bot

import "github.com/line/line-bot-sdk-go/linebot"

const helpMsg = `Usage: [@|:]mention [command] [args]

- Mention bot by: @catalyst, @cat, @tr, @th, @en, or :catalyst, :cat, :tr, :th, :en.
- Commands: 
   + translate: Translate last N messages to English.
   + weather:   Report weather info.
   + help/?:    Show this help message.
`

func (c *Catalyst) help(replyTo string) error {
	_, err := c.bot.PushMessage(replyTo, linebot.NewTextMessage(helpMsg)).Do()
	return err
}
