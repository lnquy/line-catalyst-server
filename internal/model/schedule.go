package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/robfig/cron/v3"
	qcron "github.com/lnquy/cron"
	log "github.com/sirupsen/logrus"

	"github.com/lnquy/line-catalyst-server/pkg/utils"
)

const EqualSignReplacer = "(^.!#*]"

type Schedule struct {
	Id        bson.ObjectId `json:"_" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	Cron      string        `json:"cron" bson:"cron"`
	Message   string        `json:"message" bson:"message"`
	ReplyTo   string        `json:"reply_to" bson:"reply_to"`
	IsDone    bool          `json:"is_done" bson:"is_done"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	LastRun   time.Time     `json:"last_run" bson:"last_run"`
}

func (s *Schedule) String() string {
	next := ""
	if !s.IsDone {
		cronSched, err := cron.ParseStandard(s.Cron)
		if err == nil && cronSched != nil {
			// _, offset := time.Now().In(utils.GlobalLocation).Zone()
			next = cronSched.Next(s.LastRun.In(utils.GlobalLocation)).Format(time.RFC3339)
		}
	}
	msg := strings.ReplaceAll(s.Message, EqualSignReplacer, "=")
	return fmt.Sprintf("Name: %s\nMessage: %s\nFinished: %v\nCron: %s\nSchedule: %s\nLast run: %s\nNext run: %s", s.Name, msg, s.IsDone, s.Cron, getCronDescription(s.Cron), s.LastRun.In(utils.GlobalLocation).Format(time.RFC3339), next)
}

func (s *Schedule) ShortString() string {
	next := ""
	if !s.IsDone {
		cronSched, err := cron.ParseStandard(s.Cron)
		if err == nil && cronSched != nil {
			// _, offset := time.Now().In(utils.GlobalLocation).Zone()
			next = cronSched.Next(s.LastRun.In(utils.GlobalLocation)).Format(time.RFC3339)
		}
	}
	return fmt.Sprintf("Sent by reminder: %s\nSchedule: %s\nNext run: %s", s.Name, getCronDescription(s.Cron), next)
}

func getCronDescription(expr string) string {
	exprDesc, err := qcron.NewDescriptor()
	if err != nil {
		log.Errorf("failed to init CRON expression descriptor: %s", err)
		return ""
	}
	desc, err := exprDesc.ToDescription(expr, qcron.Locale_en)
	if err != nil {
		log.Errorf("failed to get CRON description: %s", err)
		return ""
	}
	return desc
}
