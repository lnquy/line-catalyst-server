package bot

import (
	"fmt"
	"strings"

	"github.com/cosiner/flag"
	log "github.com/sirupsen/logrus"

	"github.com/lnquy/line-catalyst-server/internal/model"
)

const (
	remindAddSubCmd    = "add"
	remindGetSubCmd    = "get"
	remindListSubCmd   = "list"
	remindDeleteSubCmd = "delete"
)

type RemindAddCmd struct {
	Name    string `names:"-n, --name"`
	Cron    string `names:"-s, --schedule"`
	Message string `names:"-m, --message"`
}

func (c *Catalyst) remind(cmdArgs []string, replyTo string) error {
	if len(cmdArgs) == 0 {
		return fmt.Errorf("remind sub command must be specified")
	}

	switch strings.ToLower(strings.TrimSpace(cmdArgs[0])) {
	case remindAddSubCmd, "set", "put":
		if err := c.handleRemindAddCmd(cmdArgs[1:], replyTo); err != nil {
			return fmt.Errorf("failed to process remind add command: %w", err)
		}
	case remindGetSubCmd, "show", "view":
		if len(cmdArgs) != 2 {
			return fmt.Errorf("schedule name must be specified. E.g. @cat remind get my_schedule")
		}
		sched, err := c.scheduleRepo.Get(cmdArgs[1], replyTo)
		if err != nil {
			return fmt.Errorf("failed to get schedule name=%s: %w", cmdArgs[1], err)
		}
		c.replyTo(replyTo, fmt.Sprintf("Found schedule: %+v", sched))
	case remindListSubCmd, "get-all":
		scheds, err := c.scheduleRepo.ListAll(replyTo)
		if err != nil {
			return fmt.Errorf("failed to list all schedules: %w", err)
		}
		c.replyTo(replyTo, fmt.Sprintf("Found schedules: %+v", scheds))
	case remindDeleteSubCmd, "del", "remove", "rm":
		if len(cmdArgs) != 2 {
			return fmt.Errorf("schedule name must be specified. E.g. @cat remind delete my_schedule")
		}
		if err := c.scheduleRepo.Delete(cmdArgs[1], replyTo); err != nil {
			return fmt.Errorf("failed to delete schedule name=%s: %w", cmdArgs[1], err)
		}
		c.replyTo(replyTo, "Schedule deleted: "+cmdArgs[1])
	default:
		return fmt.Errorf("unknown sub command of remind: %s", cmdArgs[0])
	}

	return nil
}

func (c *Catalyst) handleRemindAddCmd(cmdArgs []string, replyTo string) error {
	log.Infof("remind add args: %v", cmdArgs)
	var addCmd RemindAddCmd
	if err := flag.NewFlagSet(flag.Flag{}).ParseStruct(&addCmd, cmdArgs...); err != nil {
		return fmt.Errorf("invalid add command. E.g.: @cat remind add --name my_schedule --schedule @everyday --message \"A message here\". Error: %w", err)
	}
	log.Infof("remind add: %+v", addCmd)

	sched := model.Schedule{
		Name:    addCmd.Name,
		Cron:    addCmd.Cron,
		Message: addCmd.Message,
		ReplyTo: replyTo,
		IsDone:  false,
	}
	if _, err := c.scheduleRepo.Create(&sched); err != nil {
		return fmt.Errorf("failed to save schedule: %w", err)
	}

	c.replyTo(replyTo, fmt.Sprintf("Schedule saved: %+v", sched))
	return nil
}

func (c *Catalyst) handleRemindGetCmd(cmdArgs []string, replyTo string) error {
	var addCmd RemindAddCmd
	if err := flag.NewFlagSet(flag.Flag{}).ParseStruct(&addCmd, cmdArgs...); err != nil {
		return fmt.Errorf("invalid add command. E.g.: @cat remind add --name my_schedule --schedule @everyday --message \"A message here\". Error: %w", err)
	}
	log.Infof("remind add: %+v", addCmd)

	sched := model.Schedule{
		Name:    addCmd.Name,
		Cron:    addCmd.Cron,
		Message: addCmd.Message,
		ReplyTo: replyTo,
		IsDone:  false,
	}
	if _, err := c.scheduleRepo.Create(&sched); err != nil {
		return fmt.Errorf("failed to save schedule: %w", err)
	}

	c.replyTo(replyTo, fmt.Sprintf("Schedule saved: %+v", sched))
	return nil
}
