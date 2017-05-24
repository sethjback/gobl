package manager

import (
	"github.com/robfig/cron"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/email"
	"github.com/sethjback/gobl/gobldb"
	"github.com/sethjback/gobl/gobldb/leveldb"
	"github.com/sethjback/gobl/keys"
	"github.com/sethjback/gobl/util/log"
)

var gDb gobldb.Database
var conf *config.Config
var schedules *cron.Cron
var signer keys.Signer
var verifiers map[string]keys.Verifier

// Init sets up the environement to run
func Init(c *config.Config) error {
	var err error
	gDb, err = leveldb.New(c.DB)
	if err != nil {
		return err
	}

	key, err := keys.OpenPrivateKey(c.Server.PrivateKey)
	if err != nil {
		return err
	}
	signer = keys.NewSigner(key)

	agents, err := gDb.AgentList()
	if err != nil {
		return err
	}

	verifiers = make(map[string]keys.Verifier)
	for _, a := range agents {
		akey, e := keys.DecodePublicKeyString(a.PublicKey)
		if e != nil {
			log.Errorf("manager", "error decoding public key for agent: %s", a.Name)
		} else {
			verifiers[a.ID] = keys.NewVerifier(akey)
		}

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
		log.Errorf("manager", "Could not init schedule list: %v", err)
		return err
	}

	for _, s := range ss {
		schedules.AddJob(s.String(), &ScheduledJob{s})
	}

	schedules.Start()

	log.Infof("scheduler", "Active Schedules: %v", schedules.Entries())

	return nil
}

// SendTestEmail checks the email configuration by attempting to send out a test email
func SendTestEmail() error {
	err := email.SendEmail(conf.Email, "This is a test email from gobl. Let me be the first to congratulate you on receiving this message: it means your email is configured correctly. Way to go!", "Gobl Coordinator")
	if err != nil {
		log.Errorf("manager", "could not sent test email: %v", err.Error())
		return err
	}

	return nil
}

func Shutdown() {
	schedules.Stop()
	gDb.Close()
}
