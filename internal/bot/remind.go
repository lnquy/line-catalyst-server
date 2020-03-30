package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosiner/flag"
	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"

	"github.com/lnquy/line-catalyst-server/internal/model"
	"github.com/lnquy/line-catalyst-server/pkg/utils"
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
		msg := "No reminder found. You can create new one by:\n@cat remind add --name my_sched --schedule \"@every 24h\" --message \"Trigger everyday\"\n@cat remind add -n s2 -s \"0 9 * * 1-5\" -m \"At 09:00 every day of week from Monday through Friday\"\n@cat remind add -n s3 -s \"@every 10s\" -m \"I'm flash!\""
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
		helpMsg := "Schedule a reminder to be sent.\n" +
			"Usage: @cat remind|reminder [command] [params]\n\n" +
			"Commands:\n" +
			"add|create: Create a new reminder." +
			"\n+   -n, --name: Reminder name/id." +
			"\n+   -s, --schedule: Cron scheduler format. See format below" +
			"\n+   -m, --message: Your reminder message.\n" +
			"get|view|show <reminder_name>: View a specific reminder.\n" +
			"list|get-all: View all scheduled reminders.\n" +
			"delete|remove <reminder_name>: Delete a scheduled reminder.\n\n" +
			"Examples:\n" +
			"@cat remind add --name s1 --schedule \"0 9 * * 1-5\" --message \"At 09:00 every day from Monday through Friday\"\n" +
			"@cat remind add -n s2 -s \"@every 10s\" -m \"I'm flash!\"\n" +
			"@cat remind get s2\n" +
			"@cat remind list\n" +
			"@cat remind delete s2\n" +
			"@cat remind help\n\n" +
			"Cron scheduler format: * * * * *\n" +
			"Field name   | Mandatory? | Allowed values  | Allowed special characters\n" +
			"----------   | ---------- | --------------  | --------------------------\n" +
			"Minutes      | Yes        | 0-59            | * / , -\n" +
			"Hours        | Yes        | 0-23            | * / , -\n" +
			"Day of month | Yes        | 1-31            | * / , - ?\n" +
			"Month        | Yes        | 1-12 or JAN-DEC | * / , -\n" +
			"Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?\n\n" +
			"Or using predefined functions: @yearly, @annually, @monthly, @weekly, @daily, @midnight, @hourly.\n"+
			"Or using intervals: @every <duration> (e.g.: @every 1h20m15s).\n" +
			"Details at: https://en.wikipedia.org/wiki/Cron"
		c.replyTo(replyTo, helpMsg)
	default:
		return fmt.Errorf("unknown sub command of remind: %s", cmdArgs[0])
	}

	return nil
}

func (c *Catalyst) handleRemindAddCmd(cmdArgs []string, replyTo string) error {
	log.Infof("remind add args: %v", strings.Join(cmdArgs, " | "))
	for i, _ := range cmdArgs {
		if strings.Contains(cmdArgs[i], "http") && strings.Contains(cmdArgs[i], "=") {
			cmdArgs[i] = strings.ReplaceAll(cmdArgs[i], "=", "\\=")
		}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("Recovered in handleRemindAddCmd: %s", r)
		}
	}()

	var addCmd RemindAddCmd
	if err := flag.NewFlagSet(flag.Flag{}).ParseStruct(&addCmd, cmdArgs...); err != nil {
		return fmt.Errorf("invalid add command. E.g.: @cat remind add --name my_schedule --schedule @everyday --message \"A message here\". Error: %s", err)
	}
	log.Infof("remind add: %+v", addCmd)

	cronSched, err := cron.ParseStandard(addCmd.Cron)
	if err != nil {
		return fmt.Errorf("invalid cron schedule format (%s): %s", addCmd.Cron, err)
	}
	log.Infof("cron sched: %#v", cronSched)

	now := time.Now().In(utils.GlobalLocation)
	sched := model.Schedule{
		Name:      addCmd.Name,
		Cron:      addCmd.Cron,
		Message:   addCmd.Message,
		ReplyTo:   replyTo,
		IsDone:    false,
		CreatedAt: now,
		LastRun:   now,
	}
	if created, err := c.scheduleRepo.Get(addCmd.Name, replyTo); err == nil && created != nil {
		return fmt.Errorf("reminder with the same name already existed")
	}
	if _, err := c.scheduleRepo.Create(&sched); err != nil {
		return fmt.Errorf("failed to save schedule: %s", err)
	}

	job := cron.NewWithLocation(utils.GlobalLocation)
	job.Schedule(cronSched, cron.FuncJob(func() {
		c.runReminder(&sched)
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

func (c *Catalyst) startAllScheduledReminders() error {
	scheds, err := c.scheduleRepo.ListAllScheduled()
	if err != nil {
		return fmt.Errorf("cannot list all scheduled reminders: %s", err)
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	for _, sched := range scheds {
		sched := sched

		job := cron.NewWithLocation(utils.GlobalLocation)
		// Note: All @every crons will be run at different time after restarted
		err := job.AddFunc(sched.Cron, func() {
			c.runReminder(sched)
		})
		if err != nil {
			return fmt.Errorf("failed to schedule reminder: %s\n%s", err, sched.String())
		}
		job.Start()
		c.schedMap[sched.ReplyTo+"/"+sched.Name] = job
	}
	log.Infof("all %d reminder(s) has been scheduled", len(scheds))
	return nil
}

func (c *Catalyst) runReminder(sched *model.Schedule) {
	sched.LastRun = time.Now().In(utils.GlobalLocation)

	unescaped := strings.ReplaceAll(sched.Message, "\\=", "=")
	c.replyTo(sched.ReplyTo, fmt.Sprintf("%s\n----------\n%s", unescaped, sched.ShortString()))

	_, _ = c.scheduleRepo.Update(sched)
}
