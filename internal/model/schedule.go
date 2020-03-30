package model

import (
	"fmt"
	"os"
	"time"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

var location *time.Location

func init() {
	locStr := os.Getenv("LOCATION")
	if locStr == "" {
		locStr = "Asia/Bangkok"
	}
	var err error
	location, err = time.LoadLocation(locStr)
	if err != nil {
		log.Panicf("failed to load location: %s", locStr)
	}
}

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
		cronSched, err := cron.Parse(s.Cron)
		if err == nil && cronSched != nil {
			next = cronSched.Next(s.LastRun).In(location).Format(time.RFC3339)
		}
	}
	return fmt.Sprintf("Name: %s\nMessage: %s\nFinished: %v\nSchedule: %s\nLast run: %s\nNext run: %s", s.Name, s.Message, s.IsDone, s.Cron, s.LastRun.Format(time.RFC3339), next)
}
