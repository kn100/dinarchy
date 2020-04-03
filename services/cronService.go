package services

import (
	"dinarchy/models"
	"errors"
	"fmt"
	"strconv"

	"github.com/robfig/cron/v3"
	tb "gopkg.in/tucnak/telebot.v2"
)

type CronService struct {
	crons map[string]*cron.Cron
	jobs  map[string]cron.EntryID
	TS    *tb.Bot // shouldn't need this.
}

func (c *CronService) Init() {
	c.crons = make(map[string]*cron.Cron)
	c.jobs = make(map[string]cron.EntryID)
}

func (c *CronService) AddAJob(j models.Job) (cron.EntryID, error) {
	tgidint, err := strconv.Atoi(j.TGID)
	if err != nil {
		panic(err)
	}
	cron := c.getCron(j.TGID)
	eid, err := cron.AddFunc(j.CronString, func() { c.TS.Send(&Recipient{TGID: tgidint}, j.Message, &tb.SendOptions{ParseMode: tb.ModeMarkdown}) })
	if err != nil {
		fmt.Println("errored on adding function")
		return eid, err
	}
	fmt.Println("Added function")
	c.jobs[j.Name] = eid
	return eid, nil
}

func (c *CronService) RemoveJob(tgid, name string) error {
	crond := c.getCron(tgid)
	if val, ok := c.jobs[name]; ok {
		fmt.Println("got key")
		crond.Remove(val)
		delete(c.jobs, name)
		if len(crond.Entries()) == 0 {
			fmt.Println("cleaning up cron")
			delete(c.crons, tgid)
		}
		fmt.Println("done")
		return nil
	}
	fmt.Println("couldn't find cron")
	return errors.New("That cron didn't exist")
}

func (c *CronService) getCron(tgid string) *cron.Cron {
	if val, ok := c.crons[tgid]; ok {
		return val
	}
	c.crons[tgid] = cron.New()
	c.crons[tgid].Start()
	return c.crons[tgid]
}
