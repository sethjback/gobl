package manager

import (
	"github.com/robfig/cron"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/email"
	"github.com/sethjback/gobl/gobldb"
	"github.com/sethjback/gobl/keys"
)

var gDb gobldb.Database
var schedules *cron.Cron
var conf config.Store
var signer keys.Signer
var verifiers map[string]keys.Verifier

// Init sets up the environement to run
func Init(c config.Store) error {
	var err error
	gDb, err = gobldb.Get(c)
	if err != nil {
		return err
	}

	conf = c

	//init existing schedules
	err = initCron()

	return err
}

func initCron() error {
	schedules = cron.New()
	ss, err := gDb.ScheduleList()
	if err != nil {
		return err
	}

	for _, s := range ss {
		schedules.AddJob(s.String(), &ScheduledJob{s})
	}

	schedules.Start()

	return nil
}

// SendTestEmail checks the email configuration by attempting to send out a test email
func SendTestEmail() error {
	err := email.SendEmail(conf, "This is a test email from gobl. Let me be the first to congratulate you on receiving this message: it means your email is configured correctly. Way to go!", "Gobl Coordinator")

	return err
}

func Shutdown() {
	schedules.Stop()
	gDb.Close()
}
