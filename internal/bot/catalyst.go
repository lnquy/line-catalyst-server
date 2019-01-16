package bot

import (
	"fmt"
	"net/http"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Catalyst struct {
	secret string
	token  string
	bot    *linebot.Client
}

func NewCatalyst(secret, token string) (*Catalyst, error) {
	lb, err := linebot.New(secret, token, linebot.WithHTTPClient(&http.Client{
		Transport: &http.Transport{
			MaxIdleConns:          300,
			MaxIdleConnsPerHost:   300,
			DisableKeepAlives:     false,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}))
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create Line bot")
	}

	return &Catalyst{
		secret: secret,
		token:  token,
		bot:    lb,
	}, nil
}

func (c *Catalyst) MessageHandler(w http.ResponseWriter, r *http.Request) {
	events, err := c.bot.ParseRequest(r)
	r.Body.Close()
	if err != nil {
		log.Errorf("bot: failed to parse request body: %v", err)
		errCode := http.StatusBadRequest
		if err != linebot.ErrInvalidSignature {
			errCode = http.StatusInternalServerError
		}
		http.Error(w, http.StatusText(errCode), errCode)
		return
	}
	if len(events) == 0 {
		log.Warnf("bot: event length equals to 0")
		w.Write([]byte("no command"))
		return
	}

	event := events[0] // TODO: Only handle one command for now

	switch event.Type {
	case linebot.EventTypeMessage:
		switch msg := event.Message.(type) {
		case *linebot.TextMessage:
			if _, err = c.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(msg.Text)).Do(); err != nil {
				log.Errorf("bot: failed to reply to user: %v", err)
				return
			}
		}
	default:
		w.Write([]byte(fmt.Sprintf("not supported event type: %v", event.Type)))
	}
}

func (c *Catalyst) handleTextMessage(event *linebot.Event) ([]byte, error) {
	return nil, nil
}
