package bot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/lnquy/line-catalyst-server/internal/config"
	"github.com/lnquy/line-catalyst-server/internal/model"
	"github.com/lnquy/line-catalyst-server/internal/repo"
)

const (
	translateCmd = "translate"
	weatherCmd   = "weather"
	airCmd       = "air"
	helpCmd      = "help"
)

type Catalyst struct {
	conf        config.Bot
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
		conf:        conf,
		bot:         lb,
		messageRepo: messageRepo,
	}, nil
}

func (c *Catalyst) MessageHandler(w http.ResponseWriter, r *http.Request) {
	events, err := parseRequest(r) // r.Body closed inside parseRequest
	if err != nil {
		log.Errorf("bot: failed to parse request body: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
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
			err = c.handleTextMessage(event, msg)
		default:
			log.Tracef("bot: unsupported message type")
		}
	default:
		log.Tracef("bot: unsupported event type: %v", event.Type)
	}

	if err != nil {
		log.Errorf("bot: request failed: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
}

// TODO: Resolve username
func (c *Catalyst) handleTextMessage(event *linebot.Event, msg *linebot.TextMessage) error {
	var err error
	m, isUserMessage := getMessage(event, msg)
	cmdArgs, triggered := isBotTriggered(msg.Text)

	// Normal message -> Save it replyTo database for later queries
	if !triggered {
		if _, err = c.messageRepo.Create(m, isUserMessage); err != nil {
			return errors.Wrapf(err, "failed replyTo save user message")
		}
		log.Debugf("bot: user text message saved")
		return nil
	}

	// Specify which user or group/room we should reply to
	replyTo := m.GroupID
	if isUserMessage {
		replyTo = m.UserID
	}
	// Bot mentioned -> Parse command and reply
	if len(cmdArgs) == 0 {
		cmdArgs = append(cmdArgs, translateCmd) // default translate command
	}
	cmdArgs[0] = strings.TrimSpace(strings.ToLower(cmdArgs[0]))
	switch cmdArgs[0] {
	case airCmd, "aqi":
		err = c.aqi(cmdArgs, replyTo)
	case weatherCmd:
		err = c.weather(cmdArgs, replyTo)
	case "?", helpCmd:
		err = c.help(replyTo)
	case translateCmd:
		err = c.translate(replyTo, isUserMessage, cmdArgs...)
	default:
		err = c.translate(replyTo, isUserMessage, translateCmd, strings.Join(cmdArgs[:], " "))
	}

	if err != nil {
		return errors.Wrapf(err, "failed to handle command: %v", cmdArgs)
	}
	return nil
}

// parseRequest is the same with linebot.ParseRequest.
// but we remove the validateSignature code, since we already validated it via
// the ValidateLineSignature middleware.
func parseRequest(r *http.Request) ([]*linebot.Event, error) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	request := &struct {
		Events []*linebot.Event `json:"events"`
	}{}
	if err = json.Unmarshal(body, request); err != nil {
		return nil, err
	}
	return request.Events, nil
}

func isBotTriggered(s string) ([]string, bool) {
	cmds, hit := "", false
	if !(strings.HasPrefix(s, "@") || strings.HasPrefix(s, ":")) {
		return nil, false
	}
	s = s[1:]
	if strings.HasPrefix(strings.ToLower(s), "catalyst") {
		hit = true
		cmds = string(s[8:])
		goto RETURN
	}
	if strings.HasPrefix(strings.ToLower(s), "cat") {
		hit = true
		cmds = string(s[3:])
		goto RETURN
	}
	if strings.HasPrefix(s, "tr") || strings.HasPrefix(s, "th") || strings.HasPrefix(s, "en") {
		hit = true
		cmds = string(s[2:])
		goto RETURN
	}

RETURN:
	if !hit {
		return nil, false
	}
	cmds = strings.TrimPrefix(cmds, ":")
	cmds = strings.TrimSpace(cmds)
	if len(cmds) == 0 {
		return []string{translateCmd}, true
	}
	return strings.Split(cmds, " "), true
}

func getMessage(event *linebot.Event, msg *linebot.TextMessage) (*model.Message, bool) {
	m := &model.Message{
		Timestamp: event.Timestamp,
		UserID:    event.Source.UserID,
		MessageID: msg.ID,
		Text:      msg.Text,
	}
	isUserMessage := false
	switch event.Source.Type {
	case linebot.EventSourceTypeUser:
		m.Type = model.MessageTypeUser
		isUserMessage = true
	case linebot.EventSourceTypeRoom:
		m.Type = model.MessageTypeGroup
		m.GroupID = event.Source.RoomID
	case linebot.EventSourceTypeGroup:
		m.Type = model.MessageTypeGroup
		m.GroupID = event.Source.GroupID
	default:
		return nil, false
	}
	return m, isUserMessage
}
