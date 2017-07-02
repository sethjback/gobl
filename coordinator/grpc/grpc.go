package grpc

import (
	"crypto/tls"
	"fmt"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/sethjback/gobl/certificates"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/goblerr"
	pb "github.com/sethjback/gobl/goblgrpc"
	"github.com/sethjback/gobl/model"
)

type key int

const (
	grpcConfigKey key = 0

	ErrorConfigMissing = "ConfigMissing"
	ErrorServerStart   = "ServerStartFailed"
)

type grpcConfig struct {
	listen   string
	hostCert *certificates.HostCert
	caCert   *certificates.CA
}

type server struct {
	grpcServer *grpc.Server
}

type GRPC struct {
	*grpcConfig
	server     *server
	clientConn *grpc.ClientConn
}

func New(caKey, coordKey model.Key, cs config.Store) (*GRPC, error) {
	sc := configFromStore(cs)
	if sc == nil {
		sc = &grpcConfig{listen: ":50001"}
	}

	ca, err := certificates.OpenCA(certificates.CertPEM([]byte(caKey.Certificate)))
	if err != nil {
		return nil, err
	}
	sc.caCert = ca

	coord, err := certificates.OpenHostCertificate(certificates.CertPEM([]byte(coordKey.Certificate)), certificates.CertPEM([]byte(coordKey.Key)))
	if err != nil {
		return nil, err
	}
	sc.hostCert = coord

	return &GRPC{grpcConfig: sc}, nil
}

func (g *GRPC) StartServer() error {
	creds := credentials.NewTLS(
		&tls.Config{
			ClientAuth:   tls.RequestClientCert,
			Certificates: []tls.Certificate{g.hostCert.Certificate},
			ClientCAs:    g.caCert.Pool,
		})

	lis, err := net.Listen("tcp", g.listen)
	if err != nil {
		return goblerr.New("Unable to start server", ErrorServerStart, err)
	}

	g.server = &server{}

	g.server.grpcServer = grpc.NewServer(grpc.Creds(creds))
	pb.RegisterCoordinatorServer(g.server.grpcServer, g.server)
	return g.server.grpcServer.Serve(lis)
}

func (g *GRPC) StopServer() {
	g.server.grpcServer.GracefulStop()
}

func (s *server) File(stream pb.Coordinator_FileServer) error {
	return nil
}

func (s *server) Finish(ctx context.Context, fr *pb.JobDefinition) (*pb.ReturnMessage, error) {
	fmt.Printf("Job definition: %+v\n\n", fr)
	return &pb.ReturnMessage{Message: "success", Code: "you rock"}, nil
}

func (g *GRPC) Client(a *model.Agent) (pb.AgentClient, error) {
	if g.clientConn == nil {
		creds := credentials.NewTLS(
			&tls.Config{
				ServerName:   "Coordinator",
				Certificates: []tls.Certificate{g.hostCert.Certificate},
				RootCAs:      g.caCert.Pool,
			})
		var err error
		g.clientConn, err = grpc.Dial(a.Address, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, err
		}
	}

	return pb.NewAgentClient(g.clientConn), nil
}

func SaveConfig(cs config.Store, env map[string]string) error {
	sc := &grpcConfig{}
	for k, v := range env {
		switch k {
		case "GRPC_LISTEN":
			sc.listen = v
		}
	}

	if sc.listen == "" {
		sc.listen = ":50001"
	}

	cs.Add(grpcConfigKey, sc)

	return nil
}

func configFromStore(cs config.Store) *grpcConfig {
	if sc, ok := cs.Get(grpcConfigKey); ok {
		return sc.(*grpcConfig)
	}
	return nil
}
