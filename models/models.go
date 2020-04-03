package models

import (
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/robfig/cron/v3"
)

type Job struct {
	gorm.Model
	TGID       string `gorm:"type:string;primary_key"`
	Name       string `gorm:"type:string;primary_key"`
	CronString string
	Message    string
	EntryID    cron.EntryID `gorm:"-"`
}

func (j *Job) Recipient() int {
	i, err := strconv.Atoi(j.TGID)
	if err != nil {
		panic(err)
	}
	return i
}
