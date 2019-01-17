package bot

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/lnquy/line-catalyst-server/pkg/translate"
)

const defaultMsgLimit = 5

func (c *Catalyst) translate(replyTo string, isReplyToUser bool, cmdArgs ...string) error {
	limit := defaultMsgLimit

	switch {
	case len(cmdArgs) <= 1:
	case len(cmdArgs) > 1:
		pl, err := strconv.Atoi(cmdArgs[1])
		if err != nil { // => Text
			goto DIRECT_TRANSLATE
		}

		limit = pl
		if limit <= 0 {
			limit = defaultMsgLimit
		}
		if limit > 20 {
			limit = 20
		}
		// TODO: Load message from db then translate and response
		return nil
	}

DIRECT_TRANSLATE:
	translated, err := translate.Translate("th", "en", strings.Join(cmdArgs[1:], " "))
	if err != nil {
		return errors.Wrapf(err, "failed to translate text message")
	}
	c.replyTo(replyTo, translated)
	return nil
}
