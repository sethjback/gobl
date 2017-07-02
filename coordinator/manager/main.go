package manager

import (
	"github.com/robfig/cron"
	"github.com/sethjback/gobl/certificates"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/coordinator/grpc"
	"github.com/sethjback/gobl/email"
	"github.com/sethjback/gobl/gobldb"
)

var gDb gobldb.Database
var schedules *cron.Cron
var conf config.Store
var grpcHup chan struct{}
var gs *grpc.GRPC

// Init sets up the environement to run
func Init(c config.Store) error {
	var err error
	gDb, err = gobldb.Get(c)
	if err != nil {
		return err
	}

	conf = c

	grpcHup = make(chan struct{})
	go func() {
		for {
			_, ok := <-grpcHup
			if !ok {
				//channel closed, exit
				break
			}
			resetGRPCServer()
		}
	}()

	grpcHup <- struct{}{}

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
	if gs != nil {
		gs.StopServer()
	}
	schedules.Stop()
	gDb.Close()
	close(grpcHup)
}

func resetGRPCServer() error {
	coordKey, _ := gDb.GetKey("Coordinator")
	caKey, _ := gDb.GetKey("CA")
	if caKey == nil {
		if gs != nil {
			gs.StopServer()
		}
	} else {
		if coordKey == nil {
			key, err := certificates.NewHostCertificate(*caKey, "Coordinator")
			if err != nil {
				return err
			}
			coordKey = key
			gDb.SaveKey("Coordinator", *key)
		}

		if gs != nil {
			gs.StopServer()
		}

		grpcServer, err := grpc.New(*caKey, *coordKey, conf)
		if err != nil {
			return err
		}

		go grpcServer.StartServer()
	}
	return nil
}
