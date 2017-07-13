package grpcclient

import (
	"crypto/tls"

	"github.com/sethjback/gobl/certificates"
	"github.com/sethjback/gobl/coordinator/grpcserver"
	"github.com/sethjback/gobl/goblerr"
	pb "github.com/sethjback/gobl/goblgrpc"
	"github.com/sethjback/gobl/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	ErrorConfig = "ConfigurationInvalid"
)

type Client struct {
	ClientConn *grpc.ClientConn
	Pb         pb.AgentClient
}

func New(a *model.Agent) (*Client, error) {
	client := &Client{}
	ca := grpcserver.CACert()
	if ca == nil {
		return nil, goblerr.New("CA certificate not configured", ErrorConfig, "must configure CA certificate")
	}

	host, err := certificates.OpenHostCertificate(certificates.CertPEM(a.Key.Certificate), certificates.CertPEM(a.Key.Key))
	if err != nil {
		return nil, err
	}

	creds := credentials.NewTLS(
		&tls.Config{
			ServerName:   a.Name,
			Certificates: []tls.Certificate{host.Certificate},
			RootCAs:      ca.Pool,
		})

	client.ClientConn, err = grpc.Dial(a.Address, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	client.Pb = pb.NewAgentClient(client.ClientConn)
	return client, nil
}
