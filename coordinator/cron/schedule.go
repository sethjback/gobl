package cron

import (
	"github.com/google/uuid"
	"github.com/robfig/cron"
	"github.com/sethjback/gobl/model"
)

type ScheduledJob struct {
	Schedule model.Schedule `json:"schedule"`
}

func (s *ScheduledJob) Run() {
	/*jdef, err := db.Get().GetJobDefinition(s.Schedule.JobDefinitionID)
	if err != nil {
		//TODO: log
	}

	a, err := db.Get().GetAgent(s.Schedule.AgentID)
	if err != nil {
		//TODO: log
	}

	_, err = NewJob(*jdef, a.ID)

	if err != nil {
		//TODO: log
	}*/
}

func NewSchedule(s model.Schedule) (string, error) {
	s.ID = uuid.New().String()

	_, err := cron.Parse(s.String())
	if err != nil {
		return "", err
	}

	if err = db.SaveSchedule(s); err != nil {
		return "", err
	}

	cronHup <- struct{}{}
	return s.ID, err
}

func UpdateSchedule(s model.Schedule) error {
	_, err := cron.Parse(s.String())
	if err != nil {
		return err
	}

	if err = db.SaveSchedule(s); err != nil {
		return err
	}

	cronHup <- struct{}{}
	return nil
}

func DeleteSchedule(id string) error {
	err := db.DeleteSchedule(id)
	if err != nil {
		return err
	}
	cronHup <- struct{}{}
	return nil
}

func ScheduleList() ([]model.Schedule, error) {
	return db.ScheduleList()
}

func Active() []*cron.Entry {
	active := make([]*cron.Entry, 0)
	if schedules != nil {
		active = schedules.Entries()
	}
	return active
}
