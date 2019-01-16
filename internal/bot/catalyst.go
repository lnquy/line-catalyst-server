package bot

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/lnquy/line-catalyst-server/internal/config"
	"github.com/lnquy/line-catalyst-server/internal/repo"
)

type Catalyst struct {
	bot         *linebot.Client
	messageRepo repo.MessageRepository
}

func NewCatalyst(conf config.Bot, messageRepo repo.MessageRepository) (*Catalyst, error) {
	lb, err := linebot.New(conf.Secret, conf.Token, linebot.WithHTTPClient(&http.Client{
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
		bot:         lb,
		messageRepo: messageRepo,
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
			c.handleTextMessage(event, msg)
		}
	default:
		w.Write([]byte(fmt.Sprintf("not supported event type: %v", event.Type)))
	}
}

func (c *Catalyst) handleTextMessage(event *linebot.Event, msg *linebot.TextMessage) ([]byte, error) {
	switch event.Source.Type {
	case linebot.EventSourceTypeUser:
		if _, err := c.bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Echo:"+msg.Text)).Do(); err != nil {
			log.Errorf("bot: failed to reply to user: %v", err)
			return nil, err
		}
	case linebot.EventSourceTypeGroup:
		text, ok := isBotTriggered(msg.Text)
		if !ok {
			return nil, nil // Just ignore normal messages
		}

		if _, err := c.bot.PushMessage(event.Source.GroupID, linebot.NewTextMessage("Group response: "+text)).Do(); err != nil {
			log.Errorf("bot: failed to reply to group: %v", err)
			return nil, err
		}
	default:
		log.Warnf("bot: source type not supported: %s", event.Source.Type)
		return nil, nil
	}
	return nil, nil
}

func isBotTriggered(s string) (string, bool) {
	msg := ""
	if strings.HasPrefix(s, "@Catalyst") || strings.HasPrefix(s, "@catalyst") {
		msg = string(s[9:])
		goto RETURN
	}
	if strings.HasPrefix(s, "@tr") || strings.HasPrefix(s, ":tr") ||
		strings.HasPrefix(s, "@th") || strings.HasPrefix(s, ":th") ||
		strings.HasPrefix(s, "@en") || strings.HasPrefix(s, ":en") {
		msg = string(s[3:])
		goto RETURN
	}

RETURN:
	msg = strings.TrimPrefix(msg, ":")
	msg = strings.TrimSpace(msg)
	if len(msg) == 0 {
		return "", false
	}
	return msg, true
}
