package cron

import (
	"github.com/robfig/cron"
	"github.com/sethjback/gobl/gobldb"
)

var schedules *cron.Cron
var cronHup chan struct{}
var db gobldb.Database

func Init(dbS gobldb.Database) error {
	db = dbS
	cronHup = make(chan struct{})

	go func() {
		for {
			_, ok := <-cronHup
			if !ok {
				//channel closed, exit
				return
			}
			restartCron()
		}
	}()

	cronHup <- struct{}{}
	return nil
}

func Shutdown() {
	close(cronHup)
	if schedules != nil {
		schedules.Stop()
	}
}

func restartCron() error {
	if schedules != nil {
		schedules.Stop()
	}
	ss, err := db.ScheduleList()
	if err != nil {
		return err
	}
	if len(ss) == 0 {
		schedules = nil
		return nil
	}

	schedules = cron.New()
	for _, s := range ss {
		schedules.AddJob(s.String(), &ScheduledJob{s})
	}

	schedules.Start()
	return nil
}
