package grpcserver

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/sethjback/gobl/coordinator/email"

	"github.com/sethjback/gobl/files"
	"github.com/sethjback/gobl/gobldb"
	"github.com/sethjback/gobl/goblerr"
	"github.com/sethjback/gobl/model"

	"github.com/sethjback/gobl/config"

	"github.com/sethjback/gobl/certificates"
	pb "github.com/sethjback/gobl/goblgrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	ErrorServerStart = "ServerStartFailed"
)

var listen string
var hostCert *certificates.HostCert
var caCert *certificates.CA
var s *server
var grpcHup chan struct{}
var db gobldb.Database

type server struct {
	grpcServer *grpc.Server
}

func Init(dbS gobldb.Database) error {
	db = dbS
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
	return nil
}

func Shutdown() {
	close(grpcHup)
	if s != nil {
		s.grpcServer.GracefulStop()
	}
}

func CACert() *certificates.CA {
	return caCert
}

func GRPCHup() {
	grpcHup <- struct{}{}
}

func (s *server) File(stream pb.Coordinator_FileServer) error {
	for {
		pbfile, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.ReturnMessage{Message: "success"})
		}
		if err != nil {
			return err
		}

		jf := model.JobFile{}
		if pbfile.State == pb.State_FAILED {
			jf.Error = pbfile.Message
			jf.State = model.StateFailed
		} else {
			jf.State = model.StateFinished
		}

		jf.File = convertFile(pbfile.File)
		err = db.SaveJobFile(pbfile.JobId, jf)
		if err != nil {
			// TODO: handle errors, probably bi-directional stream
			fmt.Printf("Error saving file: %+v\n", err)
		}
	}
}

func (s *server) Restore(in *pb.RestoreRequest, stream pb.Coordinator_RestoreServer) error {
	job, err := db.GetJob(in.GetId())
	if err != nil {
		return err
	}

	for _, f := range job.Definition.Files {
		stream.Send(&pb.FileRequest{
			JobId: job.ID,
			File:  convertFromFile(f),
		})
	}

	return nil
}

func (s *server) State(ctx context.Context, fr *pb.JobState) (*pb.ReturnMessage, error) {
	job, err := db.GetJob(fr.GetId())
	if err != nil {
		return &pb.ReturnMessage{Message: err.Error(), Code: "GetJobFailed"}, errors.New("Unable to find job")
	}
	job.Meta.Total = int(fr.TotalFiles)
	job.Meta.Message = fr.Message
	job.Meta.State = fr.GetState().String()
	if fr.GetState() == pb.State_FINISHED || fr.GetState() == pb.State_FAILED {
		job.Meta.End = time.Now().UTC()
	}

	db.SaveJob(*job)

	if (fr.State == pb.State_FINISHED || fr.State == pb.State_FAILED) && email.Configured() {
		err = email.StateNotification(*job.Meta, job.ID, "GOBL Job Update")
		if err != nil {
			fmt.Println(err)
		}
	}

	return &pb.ReturnMessage{Message: "success", Code: "you rock"}, nil
}

func convertFile(pbFile *pb.File) files.File {
	file := files.File{}
	file.Signature = files.Signature{
		Path:          pbFile.Signature.Path,
		Modifications: pbFile.Signature.Modifications,
		Hash:          pbFile.Signature.Hash,
	}

	file.Meta = files.Meta{
		Mode: pbFile.Meta.Mode,
		GID:  int(pbFile.Meta.Gid),
		UID:  int(pbFile.Meta.Uid),
	}

	return file
}

func convertFromFile(f files.File) *pb.File {
	pbFile := &pb.File{}
	pbFile.Signature = &pb.Signature{
		Path:          f.Signature.Path,
		Modifications: f.Signature.Modifications,
		Hash:          f.Signature.Hash,
	}

	pbFile.Meta = &pb.Meta{
		Mode: f.Mode,
		Uid:  int32(f.UID),
		Gid:  int32(f.GID),
	}

	return pbFile
}

func SaveConfig(cs config.Store, env map[string]string) error {
	for k, v := range env {
		switch k {
		case "GRPC_LISTEN":
			listen = v
		}
	}
	if listen == "" {
		listen = ":50001"
	}

	return nil
}

func resetGRPCServer() error {
	coordKey, _ := db.GetKey("Coordinator")
	caKey, _ := db.GetKey("CA")
	if caKey == nil {
		if s != nil {
			s.grpcServer.GracefulStop()
			s = nil
		}
	} else {
		if coordKey == nil {
			key, err := certificates.NewHostCertificate(*caKey, "Coordinator")
			if err != nil {
				return err
			}
			coordKey = key
			db.SaveKey("Coordinator", *key)
		}

		if s != nil {
			s.grpcServer.GracefulStop()
		}

		s = &server{}
		var err error
		caCert, err = certificates.OpenCA(certificates.CertPEM([]byte(caKey.Certificate)))
		if err != nil {
			return err
		}

		hostCert, err = certificates.OpenHostCertificate(certificates.CertPEM([]byte(coordKey.Certificate)), certificates.CertPEM([]byte(coordKey.Key)))
		if err != nil {
			return err
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

		s.grpcServer = grpc.NewServer(grpc.Creds(creds))
		pb.RegisterCoordinatorServer(s.grpcServer, s)
		go s.grpcServer.Serve(lis)
	}
	return nil
}
