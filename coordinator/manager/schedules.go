package manager

import (
	"github.com/robfig/cron"
	"github.com/sethjback/gobble/spec"
	"github.com/sethjback/gobble/util/log"
)

type ScheduledJob struct {
	Schedule *spec.Schedule `json:"schedule"`
}

func (s *ScheduledJob) Run() {
	jobID, err := RunBackup(s.Schedule.Backup)
	if err != nil {
		log.Errorf("scheduler", "Could not run scheduled Backup: %v. Error: ", *s, err)
	}
	log.Infof("scheduler", "Scheduled Job started. ID: %v", jobID)
}

func AddSchedule(s *spec.Schedule) error {

	_, err := cron.Parse(s.String())
	if err != nil {
		return err
	}

	if err = gDb.AddSchedule(s); err != nil {
		return err
	}

	schedules.Stop()
	return initCron()
}

func UpdateSchedule(s *spec.Schedule) error {
	_, err := cron.Parse(s.String())
	if err != nil {
		return err
	}

	if err = gDb.UpdateSchedule(s); err != nil {
		return err
	}

	schedules.Stop()
	return initCron()
}

func DeleteSchedule(id int) error {
	return gDb.DeleteSchedule(id)
}

func ScheduleList() ([]*spec.Schedule, error) {
	return gDb.ScheduleList()
}

func GetSchedule(id int) (*spec.Schedule, error) {
	return gDb.GetSchedule(id)
}

func CronSchedules() []*cron.Entry {
	return schedules.Entries()
}
