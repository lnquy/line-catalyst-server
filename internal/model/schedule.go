package model

import (
	"fmt"
	"time"

	"github.com/robfig/cron"

	"github.com/lnquy/line-catalyst-server/pkg/utils"
)

type Schedule struct {
	Name      string    `json:"name" bson:"_id"`
	Cron      string    `json:"cron" bson:"cron"`
	Message   string    `json:"message" bson:"message"`
	ReplyTo   string    `json:"reply_to" bson:"reply_to"`
	IsDone    bool      `json:"is_done" bson:"is_done"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	LastRun   time.Time `json:"last_run" bson:"last_run"`
}

func (s *Schedule) String() string {
	next := ""
	if !s.IsDone {
		cronSched, err := cron.ParseStandard(s.Cron)
		if err == nil && cronSched != nil {
			_, offset := time.Now().In(utils.GlobalLocation).Zone()
			next = cronSched.Next(s.LastRun).Add(time.Second * time.Duration(-offset)).In(utils.GlobalLocation).Format(time.RFC3339)
		}
	}
	return fmt.Sprintf("Name: %s\nMessage: %s\nFinished: %v\nSchedule: %s\nLast run: %s\nNext run: %s", s.Name, s.Message, s.IsDone, s.Cron, s.LastRun.In(utils.GlobalLocation).Format(time.RFC3339), next)
}

func (s *Schedule) ShortString() string {
	next := ""
	if !s.IsDone {
		cronSched, err := cron.ParseStandard(s.Cron)
		if err == nil && cronSched != nil {
			_, offset := time.Now().In(utils.GlobalLocation).Zone()
			next = cronSched.Next(s.LastRun).Add(time.Second * time.Duration(-offset)).In(utils.GlobalLocation).Format(time.RFC3339)
		}
	}
	return fmt.Sprintf("Sent by reminder: %s\nNext run: %s", s.Name, next)
}
