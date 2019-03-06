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
	"github.com/lnquy/line-catalyst-server/internal/pkg/usermap"
	"github.com/lnquy/line-catalyst-server/internal/repo"
)

const (
	translateCmd = "translate"
	weatherCmd   = "weather"
	airCmd       = "air"
	jokeCmd      = "joke"
	helpCmd      = "help"
)

type Catalyst struct {
	conf        config.Bot
	bot         *linebot.Client
	messageRepo repo.MessageRepository
	userRepo    repo.UserRepository
	um          *usermap.UserMap
}

func NewCatalyst(conf config.Bot, messageRepo repo.MessageRepository, userRepo repo.UserRepository) (*Catalyst, error) {
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

	um, err := usermap.New(userRepo)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create user map")
	}

	c := &Catalyst{
		conf:        conf,
		bot:         lb,
		messageRepo: messageRepo,
		userRepo:    userRepo,
		um:          um,
	}

	if err := c.initJokers(); err != nil {
		return nil, errors.Wrapf(err, "failed to init jokers")
	}
	return c, nil
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
	case linebot.EventTypeMessage: // User sends message
		log.Debugf("bot: hit message event")
		switch msg := event.Message.(type) {
		case *linebot.TextMessage:
			err = c.handleTextMessage(event, msg)
		default:
			log.Tracef("bot: unsupported message type")
		}
	case linebot.EventTypeFollow: // User follows/unblocks the bot
		log.Debugf("bot: hit follow event")
		c.replyTo(event.Source.UserID, greetingMsg)
		c.resolveUsername(event.Source.UserID, "", "")
	case linebot.EventTypeMemberJoined: //  User(s) just join a group/room bot already in
		log.Debugf("bot: hit member joined event")
		for _, mem := range event.Members {
			if event.Source.Type == linebot.EventSourceTypeGroup {
				c.resolveUsername(mem.UserID, event.Source.GroupID, "")
			} else if event.Source.Type == linebot.EventSourceTypeRoom {
				c.resolveUsername(mem.UserID, "", event.Source.RoomID)
			}
		}
	case linebot.EventTypeJoin: // Bot joins a group/room
		log.Debugf("bot: hit group/room join event")
		c.handleJoinEvent(event)
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

func (c *Catalyst) handleTextMessage(event *linebot.Event, msg *linebot.TextMessage) error {
	var err error
	m, isUserMessage := getMessage(event, msg)
	cmdArgs, triggered := isBotTriggered(msg.Text)

	// Normal message -> Save it replyTo database for later queries
	if !triggered {
		m.Username = c.resolveUsername(m.UserID, m.GroupID, m.GroupID)
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
	case jokeCmd, "fun":
		err = c.joke(cmdArgs, replyTo)
	case "?", helpCmd:
		err = c.help(replyTo)
	case translateCmd:
		err = c.translate(replyTo, isUserMessage, cmdArgs...)
	default:
		err = c.translate(replyTo, isUserMessage, translateCmd, strings.Join(cmdArgs[:], " "))
	}

	if err != nil {
		c.replyTo(replyTo, "Sorry. Something wrong happened :(")
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
	if strings.HasPrefix(s, "en2th") {
		hit = true
		cmds = string(s[5:])
		c.conf.Translation.SourceLang = "en"
		c.conf.Translation.TargetLang = "th"
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

func (c *Catalyst) resolveUsername(uid, gid, rid string) string {
	username := c.um.Get(uid)
	if username != "" {
		return username
	}

	// Try to query username from group/room first,
	// if failed then try to query directly from their profile.
	if gid != "" {
		up, err := c.bot.GetGroupMemberProfile(gid, uid).Do()
		if err == nil {
			c.updateUserMapInfo(up)
			return up.DisplayName
		}
	}
	if rid != "" {
		up, err := c.bot.GetRoomMemberProfile(rid, uid).Do()
		if err == nil {
			c.updateUserMapInfo(up)
			return up.DisplayName
		}
	}
	up, err := c.bot.GetProfile(uid).Do()
	if err != nil {
		// Poison the cache so always use uid as display name for this user
		c.um.Set(uid, uid)
		return uid
	}
	c.updateUserMapInfo(up)
	return up.DisplayName
}

func (c *Catalyst) updateUserMapInfo(up *linebot.UserProfileResponse) {
	if _, err := c.userRepo.Create(&model.User{
		ID:            up.UserID,
		Name:          up.DisplayName,
		PictureURL:    up.PictureURL,
		StatusMessage: up.StatusMessage,
	}); err != nil {
		log.Errorf("bot: failed to save user info to database, cache can be missed later: %v", err) // Log only
	}
	c.um.Set(up.UserID, up.DisplayName)
}

// Note: GetGroupMemberIDs and GetRoomMemberIDs APIs requires Line@ approved account.
// https://developers.line.biz/en/reference/messaging-api/#get-group-member-user-ids
func (c *Catalyst) handleJoinEvent(event *linebot.Event) {
	switch event.Source.Type {
	case linebot.EventSourceTypeGroup:
		c.replyTo(event.Source.GroupID, greetingMsg)
		nextToken := ""
		gid := event.Source.GroupID
		for {
			mems, err := c.bot.GetGroupMemberIDs(gid, nextToken).Do()
			if err != nil {
				log.Errorf("bot: failed to get list of userIDs in the %s group: %v", gid, err)
				return
			}
			for _, uid := range mems.MemberIDs {
				c.resolveUsername(uid, gid, "")
			}
			nextToken = mems.Next
			if nextToken == "" {
				return
			}
		}
	case linebot.EventSourceTypeRoom:
		c.replyTo(event.Source.RoomID, greetingMsg)
		nextToken := ""
		rid := event.Source.RoomID
		for {
			mems, err := c.bot.GetRoomMemberIDs(rid, nextToken).Do()
			if err != nil {
				log.Errorf("bot: failed to get list of userIDs in the %s room: %v", rid, err)
				return
			}
			for _, uid := range mems.MemberIDs {
				c.resolveUsername(uid, "", rid)
			}
			nextToken = mems.Next
			if nextToken == "" {
				return
			}
		}
	}
}
