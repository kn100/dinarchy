package services

import (
	"dinarchy/models"
	"dinarchy/utils"
	"errors"

	"github.com/jinzhu/gorm"
)

type JobService struct {
	DB *gorm.DB
	CS CronService
}

func (js *JobService) AddJob(tgid string, cr utils.CreateRequest) error {
	// TODO: Verify that the job doesn't exist since gorm doesn't do primary keys with sqlite3. Either that or switch to a real db.
	job := models.Job{TGID: tgid, Name: cr.Name, CronString: cr.Cronstring, Message: cr.Message}
	errs := js.DB.Create(&job).GetErrors()
	if len(errs) != 0 { // TODO: Handle better
		return errs[0]
	}
	_, err := js.CS.AddAJob(job)
	if err != nil {
		js.RemoveJob(tgid, cr.Name)
		return errors.New("couldn't create the cron job for some reason. Your cron string was probably not shaped correctly")
	}
	return nil
}

func (js *JobService) RemoveJob(tgid, name string) error {
	err := js.CS.RemoveJob(tgid, name)
	if err != nil {
		return err
	}
	errs := js.DB.Where("tg_id = ? AND name = ?", tgid, name).Unscoped().Delete(models.Job{}).GetErrors()
	if len(errs) != 0 { // TODO: Handle better
		return errs[0]
	}
	return nil
}

func (js *JobService) GetJobs(tgid string) []models.Job {
	var jobs []models.Job
	js.DB.Where("tg_id = ?", tgid).Find(&jobs)
	return jobs
}

func (js *JobService) LoadJobs() error {
	var jobs []models.Job
	errs := js.DB.Find(&jobs).GetErrors()
	if len(errs) != 0 { // TODO: Handle better
		return errs[0]
	}
	for _, j := range jobs {
		_, err := js.CS.AddAJob(j)
		if err != nil {
			panic("AAAAAA")
		}
	}
	return nil
}
