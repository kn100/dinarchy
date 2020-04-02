package dcron

import (
	"errors"

	"github.com/davecgh/go-spew/spew"
	"github.com/robfig/cron/v3"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Dcron struct {
	Cron   *cron.Cron
	Jobs   map[string]DcronJob
	Tguser *tb.User
}

type DcronJob struct {
	Name       string
	CronString string
	Message    string
	EntryID    cron.EntryID
}

func (d *Dcron) AddJob(cronstring, name, message string, fun func(tb.Recipient, interface{}, ...interface{}) (*tb.Message, error)) (string, error) {
	if d.Jobs == nil {
		spew.Dump("initialising map")
		d.Jobs = make(map[string]DcronJob)
	}
	entryID, err := d.Cron.AddFunc(cronstring, func() { fun(tb.Recipient(d.Tguser), message) })
	if err != nil {
		return "", err
	}
	d.Jobs[name] = DcronJob{Name: name, CronString: cronstring, Message: message, EntryID: entryID}
	spew.Dump(d.Jobs)
	return name, nil
}

func (d *Dcron) RemoveJob(name string) error {
	spew.Dump(d.Jobs)
	spew.Dump(name)
	if _, ok := d.Jobs[name]; !ok {
		return errors.New("that job does not exist")
	}
	d.Cron.Remove(d.Jobs[name].EntryID)
	delete(d.Jobs, name)
	return nil
}
