package grpcserver

import (
	"crypto/tls"
	"errors"
	"net"

	"github.com/sethjback/gobl/config"

	"github.com/sethjback/gobl/agent/job"
	"github.com/sethjback/gobl/certificates"
	"github.com/sethjback/gobl/engine"
	"github.com/sethjback/gobl/goblerr"
	pb "github.com/sethjback/gobl/goblgrpc"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/modification"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	ErrorServerStart = "StartServerFailed"
	ErrorCreateJob   = "CreateJobFailed"
	ErrorCancelJob   = "CancelJobFailed"
)

type server struct {
	grpcServer *grpc.Server
}

var listen string
var hostCert *certificates.HostCert
var caCert *certificates.CA
var s *server

func Init() error {
	if listen == "" || hostCert == nil || caCert == nil {
		return goblerr.New("Unable to start server: config invalid", ErrorServerStart, nil)
	}
	creds := credentials.NewTLS(
		&tls.Config{
			ClientAuth:   tls.RequestClientCert,
			Certificates: []tls.Certificate{hostCert.Certificate},
			ClientCAs:    caCert.Pool,
		})

	lis, err := net.Listen("tcp", listen)
	if err != nil {
		return goblerr.New("Unable to start server", ErrorServerStart, err)
	}

	s = &server{}

	s.grpcServer = grpc.NewServer(grpc.Creds(creds))
	pb.RegisterAgentServer(s.grpcServer, s)
	go s.grpcServer.Serve(lis)
	return nil
}

func Shutdown() {
	s.grpcServer.GracefulStop()
}

func SaveConfig(cs config.Store, env map[string]string) error {
	hcPath := ""
	hkPath := ""
	caPath := ""
	for k, v := range env {
		switch k {
		case "GRPC_LISTEN":
			listen = v
		case "HOST_CERT":
			hcPath = v
		case "HOST_KEY":
			hkPath = v
		case "CA_CERT":
			caPath = v
		}
	}

	if hcPath == "" || hkPath == "" {
		return errors.New("Must provide host certificate and key")
	}
	if caPath == "" {
		return errors.New("Must provide CA certificate")
	}

	var err error
	hostCert, err = certificates.OpenHostCertificate(certificates.CertPath(hcPath), certificates.CertPath(hkPath))
	if err != nil {
		return err
	}

	caCert, err = certificates.OpenCA(certificates.CertPath(caPath))
	if err != nil {
		return err
	}

	if listen == "" {
		listen = ":50001"
	}

	return nil
}

func (s *server) Backup(ctx context.Context, fr *pb.JobDefinition) (*pb.ReturnMessage, error) {
	def := buildDefinition(fr)
	j := model.Job{
		ID:         fr.Id,
		Definition: def,
	}

	err := job.NewBackup(j)
	if err != nil {
		return &pb.ReturnMessage{Message: "unable to create job: " + err.Error(), Code: ErrorCreateJob}, nil
	}

	return &pb.ReturnMessage{Message: "success"}, nil
}

func (s *server) Restore(ctx context.Context, fr *pb.JobDefinition) (*pb.ReturnMessage, error) {
	def := buildDefinition(fr)
	j := model.Job{
		ID:         fr.Id,
		Definition: def,
	}
	err := job.NewRestore(j)
	if err != nil {
		return &pb.ReturnMessage{Message: "unable to create job: " + err.Error(), Code: ErrorCreateJob}, nil
	}

	return &pb.ReturnMessage{Message: "success"}, nil
}

func (s *server) Cancel(ctx context.Context, fr *pb.JobDefinition) (*pb.ReturnMessage, error) {
	rm := &pb.ReturnMessage{}
	err := job.Cancel(fr.Id)

	if err != nil {
		rm.Message = "cancel job failed: " + err.Error()
		rm.Code = ErrorCancelJob
	} else {
		rm.Message = "success"
	}
	return rm, nil
}

func (s *server) State(ctx context.Context, fr *pb.StateRequest) (*pb.JobState, error) {
	js := &pb.JobState{}
	meta, err := job.State(fr.Id)
	if err != nil {
		gerr := err.(goblerr.Error)
		js.Message = err.Error()
		if gerr.Code == job.ErrorFindJob {
			js.State = pb.State_NOTFOUND
		}
	} else {
		js.Message = meta.Message
		js.TotalFiles = int32(meta.Total)
		js.CompletedFiles = int32(meta.Complete)
		switch meta.State {
		case model.StateRunning, model.StateNotification:
			js.State = pb.State_RUNNING
		case model.StateCanceled, model.StateCanceling:
			js.State = pb.State_CANCELED
		case model.StateFailed:
			js.State = pb.State_FAILED
		}
	}
	return js, nil
}

func buildDefinition(fr *pb.JobDefinition) *model.JobDefinition {
	def := &model.JobDefinition{ID: fr.Id}
	for _, v := range fr.To {
		t := engine.Definition{}
		t.Name = v.Name
		t.Options = make(map[string]string)
		for _, val := range v.Options {
			t.Options[val.Name] = val.Value
		}
		def.To = append(def.To, t)
	}

	if fr.From != nil {
		t := &engine.Definition{}
		t.Name = fr.From.Name
		t.Options = make(map[string]string)
		for _, val := range fr.From.Options {
			t.Options[val.Name] = val.Value
		}
		def.From = t
	}

	for _, v := range fr.Modifications {
		t := modification.Definition{}
		t.Name = v.Name
		t.Options = make(map[string]string)
		for _, val := range v.Options {
			t.Options[val.Name] = val.Value
		}
		def.Modifications = append(def.Modifications, t)
	}

	for _, v := range fr.Paths {
		t := model.Path{}
		t.Root = v.Root
		t.Excludes = append(t.Excludes, v.Excludes...)
		def.Paths = append(def.Paths, t)
	}

	return def
}
