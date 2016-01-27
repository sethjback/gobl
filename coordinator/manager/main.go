package manager

import (
	"crypto/rsa"
	"fmt"
	"strings"

	"github.com/robfig/cron"
	"github.com/sethjback/gobl/coordinator/gobldb"
	"github.com/sethjback/gobl/coordinator/gobldb/sqlite"
	"github.com/sethjback/gobl/keys"
	"github.com/sethjback/gobl/util/log"
)

var gDb gobldb.Database
var keyManager *keys.Manager
var hostConfig map[string]interface{}
var schedules *cron.Cron

// Init sets up the environement to run
func Init(config map[string]interface{}) error {
	gDb = &sqlite.SQLite{}
	err := gDb.Init(config)

	//err = db.Init()
	if err != nil {
		return err
	}

	keyManager = &keys.Manager{PublicKeys: make(map[string]*rsa.PublicKey)}
	key, err := keys.OpenPrivateKey(config["PrivateKeyPath"].(string))
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

	hostConfig = config

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
