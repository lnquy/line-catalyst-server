package bot

import "github.com/line/line-bot-sdk-go/linebot"

const (
	helpMsg = `Usage: [@|:]mention [command] [args]

- Mention bot by: @catalyst, @cat, @bot, or :catalyst, :cat, :bot.

- Commands: 
   + translate: Translate last N messages from auto detect language to English.
   + th2en: Translate last N messages from Thai to English.
   + en2th: Translate last N messages from English to Thai.
   + th2vi: Translate last N messages from Thai to Vietnamese.
   + vi2th: Translate last N messages from Vietnamese to Thai.
   + weather: Report current weather info.
   + air/aqi: Report Air Quality Index (AQI).
   + joke/fun: Tell me a joke.
   + help/?: Show this help message.

- Examples:
   + @catalyst               // Translate last 5 messages to English.
   + @cat ผมรักคุณ            // Translate "ผมรักคุณ" to English.
   + @cat translate 15       // Translate last 15 messages to English.
   + @bot en2th              // Translate last 5 messages from English to Thai.
   + @bot en2th Hello there! // Translate "Hello there!" to Thai.
   + :cat weather            // Report current weather.
   + @cat weather hcmc       // Report current weather for HoChiMinh city.
   + @cat air                // Report current Air Quality Index (AQI).
   + @cat aqi hanoi          // Report current AQI for Hanoi city.
   + @cat joke               // Tell a joke.

Source code, help or report issue at: http://git.io/fhBYS. 
`
	greetingMsg = `Hello,
I'm Catalyst :).

Just mention me in chat box as @cat, @bot, @catalyst, or :cat, :bot, :catalyst; then I will do some boring stuffs for you.
Type "@cat ?" to display help message.

My source code can be found at: https://git.io/fhBYS.
`
)

func (c *Catalyst) help(replyTo string) error {
	_, err := c.bot.PushMessage(replyTo, linebot.NewTextMessage(helpMsg)).Do()
	return err
}
