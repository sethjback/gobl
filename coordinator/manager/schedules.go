package manager

import (
	"github.com/google/uuid"
	"github.com/robfig/cron"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/util/log"
)

// Implements
type ScheduledJob struct {
	Schedule model.Schedule `json:"schedule"`
}

func (s *ScheduledJob) Run() {
	jdef, err := gDb.GetJobDefinition(s.Schedule.JobDefinitionID)
	if err != nil {
		log.Errorf("scheduler", "Could not job definition for schedule %+v. Error: %v", *s, err)
	}

	a, err := gDb.GetAgent(s.Schedule.AgentID)
	if err != nil {
		log.Errorf("scheduler", "Could not find agent for schedule %+v. Error: %v", *s, err)
	}

	jID, err := NewJob(*jdef, a.ID)

	if err != nil {
		log.Errorf("scheduler", "Could not run scheduled job: %v. Error: %v", *s, err)
	}

	log.Infof("scheduler", "Scheduled Job started. ID: %v", jID)
}

func NewSchedule(s model.Schedule) (string, error) {
	s.ID = uuid.New().String()

	_, err := cron.Parse(s.String())
	if err != nil {
		return "", err
	}

	if err = gDb.SaveSchedule(s); err != nil {
		return "", err
	}

	schedules.Stop()
	err = initCron()
	return s.ID, err
}

func UpdateSchedule(s model.Schedule) error {
	_, err := cron.Parse(s.String())
	if err != nil {
		return err
	}

	if err = gDb.SaveSchedule(s); err != nil {
		return err
	}

	schedules.Stop()
	return initCron()
}

func DeleteSchedule(id string) error {
	return gDb.DeleteSchedule(id)
}

func ScheduleList() ([]model.Schedule, error) {
	return gDb.ScheduleList()
}

func CronSchedules() []*cron.Entry {
	return schedules.Entries()
}
