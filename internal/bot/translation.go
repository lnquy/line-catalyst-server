package bot

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/lnquy/line-catalyst-server/pkg/translate"
)

const defaultMsgLimit = 5

func (c *Catalyst) translate(cmdArgs []string, replyTo string, isReplyToUser bool) error {
	if len(cmdArgs) > 2 { // Translate the text after the command
		translated, err := translate.Translate("th", "en", strings.Join(cmdArgs[1:], " "))
		if err != nil {
			return errors.Wrapf(err, "failed to translate text message")
		}
		c.replyTo(replyTo, "TH -> EN:\n" + translated)
	}

	limit := defaultMsgLimit
	if len(cmdArgs) == 2 {
		pl, err := strconv.Atoi(cmdArgs[1])
		if err == nil {
			limit = pl
		}
	}
	if limit <= 0 {
		limit = defaultMsgLimit
	}
	if limit > 20 {
		limit = 20
	}

	// TODO: Load message from db then translate and response
	return nil
}
