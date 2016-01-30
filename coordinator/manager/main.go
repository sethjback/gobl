package manager

import (
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/robfig/cron"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/coordinator/gobldb"
	"github.com/sethjback/gobl/coordinator/gobldb/sqlite"
	"github.com/sethjback/gobl/email"
	"github.com/sethjback/gobl/keys"
	"github.com/sethjback/gobl/util/log"
)

var gDb gobldb.Database
var keyManager *keys.Manager
var conf *config.Config
var schedules *cron.Cron

// Init sets up the environement to run
func Init(c *config.Config) error {
	gDb = &sqlite.SQLite{}
	err := gDb.Init(c.DB)

	//err = db.Init()
	if err != nil {
		return err
	}

	keyManager = &keys.Manager{PublicKeys: make(map[string]*rsa.PublicKey)}
	key, err := keys.OpenPrivateKey(c.Host.PrivateKeyPath)
	if err != nil {
		return err
	}

	keyManager.PrivateKey = key

	agents, err := gDb.AgentList()
	if err != nil {
		return err
	}

	for _, a := range agents {
		key, err := keys.DecodePublicKeyString(a.PublicKey)
		if err != nil {
			fmt.Println("Error decoding key for agent: ", a)
		} else {
			ip := strings.Split(a.Address, ":")
			keyManager.PublicKeys[ip[0]] = key
		}

	}

	conf = c

	//init existing schedules
	if err = initCron(); err != nil {
		return err
	}

	return nil
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

// VerifySignature checks the incoming request agains
func VerifySignature(host string, signed []byte, signature string) error {
	ip := strings.Split(host, ":")
	key, err := keyManager.KeyForHost(ip[0])
	if err != nil {
		return err
	}

	return keys.VerifySignature(key, signed, signature)
}

// SendTestEmail checks the email configuration by attempting to send out a test email
func SendTestEmail() error {
	err := email.SendEmail(conf.Email, "This is a test email from gobl. Let me be the first to congratulate you on receiving this message: it means your email is configured correctly. Way to go!", "Netfung Gobl Coordinator")
	if err != nil {
		log.Errorf("manager", "could not sent test email: %v", err.Error())
		return err
	}

	return nil
}
