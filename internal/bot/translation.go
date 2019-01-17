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
		messages, err := c.messageRepo.ListLastMessages(replyTo, limit, isReplyToUser)
		if err != nil {
			return errors.Wrapf(err, "failed to load messages for %s", replyTo)
		}

		for _, m := range messages {
			text += getTextFromMessage(m) + "\n-----\n"
		}
		text = strings.TrimSuffix(text, "\n-----\n")
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
		translated = "Sorry. No text translated :("
	}
	c.replyTo(replyTo, translated)
	return nil
}

// TODO: Username
func getTextFromMessage(m *model.Message) string {
	return fmt.Sprintf("%s (%s): %s", m.UserID, m.Timestamp.Format("01/02/2006 15:04"), m.Text)
}
