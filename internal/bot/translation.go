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
	var messages []*model.Message

	switch {
	case len(cmdArgs) <= 1:
		var err error
		messages, text, err = c.prepareMessageText(replyTo, limit, isReplyToUser)
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
		messages, text, err = c.prepareMessageText(replyTo, limit, isReplyToUser)
		if err != nil {
			return errors.Wrapf(err, "failed to get user messages to translate")
		}
		goto TRANSLATE
	}

TRANSLATE:
	if text == "" {
		c.replyTo(replyTo, "Sorry. No message to translate :(")
		return nil
	}
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

	// Replace name-time placeholder by actual value
	for i, m := range messages {
		translated = strings.Replace(
			translated,
			fmt.Sprintf("${{%d}}", i),
			fmt.Sprintf("%s (%s)", m.Username, m.Timestamp.Format("02/01/2006 15:04")),
			1,
		)
	}
	c.replyTo(replyTo, translated)
	return nil
}

func (c *Catalyst) prepareMessageText(replyTo string, limit int, isUserMessage bool) ([]*model.Message, string, error) {
	messages, err := c.messageRepo.ListLastMessages(replyTo, limit, isUserMessage)
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to load messages for %s", replyTo)
	}

	text := ""
	for i, m := range messages {
		text += fmt.Sprintf("${{%d}}: %s\n-----\n", i, m.Text)
	}
	return messages, strings.TrimSuffix(text, "\n-----\n"), nil
}
