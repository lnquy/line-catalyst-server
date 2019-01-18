package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/lnquy/line-catalyst-server/internal/model"
	"github.com/lnquy/line-catalyst-server/pkg/translate"
)

const defaultMsgLimit = 5

func (c *Catalyst) translate(replyTo string, isReplyToUser bool, cmdArgs ...string) error {
	limit := defaultMsgLimit
	text := ""

	switch {
	case len(cmdArgs) <= 1:
		var err error
		text, err = c.prepareMessageText(replyTo, limit, isReplyToUser)
		if err != nil {
			return errors.Wrapf(err, "failed to get user messages to translate")
		}
		goto TRANSLATE
	case len(cmdArgs) > 1:
		pl, err := strconv.Atoi(cmdArgs[1])
		if err != nil { // => Text
			text = strings.Join(cmdArgs[1:], " ")
			goto TRANSLATE
		}

		limit = pl
		if limit <= 0 {
			limit = defaultMsgLimit
		}
		if limit > 20 {
			limit = 20
		}
		text, err = c.prepareMessageText(replyTo, limit, isReplyToUser)
		if err != nil {
			return errors.Wrapf(err, "failed to get user messages to translate")
		}
		goto TRANSLATE
	}

TRANSLATE:
	translated, err := translate.GoogleTranslate(
		c.conf.Translation.SourceLang,
		c.conf.Translation.TargetLang,
		text,
	)
	if err != nil {
		return errors.Wrapf(err, "failed to translate text message")
	}
	if translated == "" {
		translated = "Sorry. No message translated :("
	}
	c.replyTo(replyTo, translated)
	return nil
}

func (c *Catalyst) prepareMessageText(replyTo string, limit int, isUserMessage bool) (string, error) {
	messages, err := c.messageRepo.ListLastMessages(replyTo, limit, isUserMessage)
	if err != nil {
		return "", errors.Wrapf(err, "failed to load messages for %s", replyTo)
	}

	text := ""
	for _, m := range messages {
		text += getTextFromMessage(m) + "\n-----\n"
	}
	return strings.TrimSuffix(text, "\n-----\n"), nil
}

// TODO: Username
func getTextFromMessage(m *model.Message) string {
	return fmt.Sprintf("%s (%s): %s", m.Username, m.Timestamp.Format("01/02/2006 15:04"), m.Text)
}
