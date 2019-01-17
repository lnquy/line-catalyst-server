package bot

import "github.com/line/line-bot-sdk-go/linebot"

const helpMsg = `Usage: [@|:]mention [command] [args]

- Mention bot by: @catalyst, @cat, @tr, @th, @en, or :catalyst, :cat, :tr, :th, :en.

- Commands: 
   + translate: Translate last N messages to English.
   + weather:   Report current weather info.
   + air/aqi:   Report Air Quality Index (AQI).
   + help/?:    Show this help message.

- Examples:
   + @catalyst           // Translate last 5 messages.
   + @cat ผมรักคุณ        // Translate "ผมรักคุณ" to English.
   + @cat translate 15   // Translate last 15 messages.
   + :cat weather        // Report current weather.
   + @cat weather hanoi  // Report current weather for Hanoi city.
   + @cat air            // Report current Air Quality Index (AQI).
   + @cat aqi hanoi      // Report current AQI for Hanoi city.
`

func (c *Catalyst) help(replyTo string) error {
	_, err := c.bot.PushMessage(replyTo, linebot.NewTextMessage(helpMsg)).Do()
	return err
}
