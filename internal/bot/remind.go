package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosiner/flag"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"

	"github.com/lnquy/line-catalyst-server/internal/model"
)

const (
	remindAddSubCmd    = "add"
	remindGetSubCmd    = "get"
	remindListSubCmd   = "list"
	remindDeleteSubCmd = "delete"
	remindHelpSubCmd   = "help"
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
	case remindAddSubCmd, "create":
		if err := c.handleRemindAddCmd(cmdArgs, replyTo); err != nil {
			return fmt.Errorf("failed to process remind add command: %s", err)
		}
	case remindGetSubCmd, "show", "view":
		if len(cmdArgs) != 2 {
			return fmt.Errorf("schedule name must be specified. E.g. @cat remind get my_schedule")
		}
		sched, err := c.scheduleRepo.Get(cmdArgs[1], replyTo)
		if err != nil {
			return fmt.Errorf("failed to get schedule name=%s: %s", cmdArgs[1], err)
		}
		c.replyTo(replyTo, sched.String())
	case remindListSubCmd, "get-all":
		scheds, err := c.scheduleRepo.ListAll(replyTo)
		if err != nil {
			return fmt.Errorf("failed to list all schedules: %s", err)
		}
		msg := `No reminder found. You can create new one by:\n@cat remind add --name my_sched --schedule "@every 24h" --message "Trigger everyday"\n@cat remind add -n s2 -s "0 9 * * 1-5" -m "At 09:00 every day of week from Monday through Friday"\n@cat remind add -n s3 -s "@every 10s" -m "I'm flash!"`
		if len(scheds) != 0 {
			msg = fmt.Sprintf("Found %d reminder(s)\n----------\n", len(scheds))
			for _, sched := range scheds {
				msg += fmt.Sprintf("%s\n-----\n", sched.String())
			}
			msg = strings.TrimSuffix(msg, "\n-----\n")
		}
		c.replyTo(replyTo, msg)
	case remindDeleteSubCmd, "remove":
		if err := c.handleRemindDeleteCmd(cmdArgs, replyTo); err != nil {
			return fmt.Errorf("failed to process remind delete command: %s", err)
		}
	case remindHelpSubCmd, "h", "?":
		helpMsg := `Help` // TODO
		c.replyTo(replyTo, helpMsg)
	default:
		return fmt.Errorf("unknown sub command of remind: %s", cmdArgs[0])
	}

	return nil
}

func (c *Catalyst) handleRemindAddCmd(cmdArgs []string, replyTo string) error {
	log.Infof("remind add args: %v", cmdArgs)
	var addCmd RemindAddCmd
	if err := flag.NewFlagSet(flag.Flag{}).ParseStruct(&addCmd, cmdArgs...); err != nil {
		return fmt.Errorf("invalid add command. E.g.: @cat remind add --name my_schedule --schedule @everyday --message \"A message here\". Error: %s", err)
	}
	log.Infof("remind add: %+v", addCmd)

	cronSched, err := cron.Parse(addCmd.Cron)
	if err != nil {
		return fmt.Errorf("invalid cron schedule format (%s): %s", addCmd.Cron, err)
	}

	now := time.Now()
	sched := model.Schedule{
		Name:      addCmd.Name,
		Cron:      addCmd.Cron,
		Message:   addCmd.Message,
		ReplyTo:   replyTo,
		IsDone:    false,
		CreatedAt: now,
		LastRun:   now,
	}
	if _, err := c.scheduleRepo.Create(&sched); err != nil {
		return fmt.Errorf("failed to save schedule: %s", err)
	}

	job := cron.New()
	job.Schedule(cronSched, cron.FuncJob(func() {
		c.replyTo(sched.ReplyTo, fmt.Sprintf("[%s] %s", sched.Name, sched.Message))

		sched.LastRun = time.Now()
		_, _ = c.scheduleRepo.Update(&sched)
	}))
	job.Start()

	c.lock.Lock()
	c.schedMap[replyTo+"/"+sched.Name] = job
	c.lock.Unlock()

	c.replyTo(replyTo, fmt.Sprintf("Reminder has been scheduled\n----------\n%s\n\nYou can manage it by: @cat remind get|delete %s", sched.String(), sched.Name))
	return nil
}

func (c *Catalyst) handleRemindDeleteCmd(cmdArgs []string, replyTo string) error {
	if len(cmdArgs) != 2 {
		return fmt.Errorf("schedule name must be specified. E.g. @cat remind delete my_schedule")
	}
	sched, err := c.scheduleRepo.Get(cmdArgs[1], replyTo)
	if err != nil {
		return fmt.Errorf("failed to get schedule: %s", err)
	}
	sched.IsDone = true

	if err := c.scheduleRepo.Delete(cmdArgs[1], replyTo); err != nil {
		return fmt.Errorf("failed to delete schedule name=%s: %s", cmdArgs[1], err)
	}

	id := replyTo + "/" + cmdArgs[1]
	var job *cron.Cron
	c.lock.Lock()
	job = c.schedMap[id]
	delete(c.schedMap, id)
	c.lock.Unlock()

	job.Stop()

	c.replyTo(replyTo, fmt.Sprintf("Reminder deleted\n----------\n%s", sched.String()))
	return nil
}
