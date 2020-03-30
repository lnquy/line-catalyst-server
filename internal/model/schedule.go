package model

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

type Schedule struct {
	Name    string `json:"name" bson:"_id"`
	Cron    string `json:"cron" bson:"cron"`
	Message string `json:"message" bson:"message"`
	ReplyTo string `json:"reply_to" bson:"reply_to"`
	IsDone  bool   `json:"is_done" bson:"is_done"`
}

func (s *Schedule) String() string {
	next := ""
	if !s.IsDone {
		cronSched, err := cron.Parse(s.Cron)
		if err == nil && cronSched != nil {
			next = cronSched.Next(time.Now()).Format(time.RFC3339)
		}
	}
	return fmt.Sprintf("Name: %s\nMessage: %s\nSchedule: %s\nNext run: %s\nIs finished: %v", s.Name, s.Message, s.Cron, next, s.IsDone)
}
