package grpc

import (
	"crypto/tls"
	"errors"
	"net"

	"github.com/sethjback/gobl/agent/coordinator"
	"github.com/sethjback/gobl/config"

	"golang.org/x/net/context"

	"github.com/sethjback/gobl/certificates"
	"github.com/sethjback/gobl/goblerr"
	pb "github.com/sethjback/gobl/goblgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

func New(cs config.Store) (*GRPC, error) {
	sc := configFromStore(cs)
	if sc == nil {
		return nil, goblerr.New("grpc server config missing", ErrorConfigMissing, nil)
	}

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
	pb.RegisterAgentServer(g.server.grpcServer, g.server)
	return g.server.grpcServer.Serve(lis)
}

func (g *GRPC) StopServer() {
	g.server.grpcServer.GracefulStop()
}

func (s *server) Backup(ctx context.Context, fr *pb.JobDefinition) (*pb.ReturnMessage, error) {
	return nil, nil
}

func (s *server) Restore(ctx context.Context, fr *pb.JobDefinition) (*pb.ReturnMessage, error) {
	return nil, nil
}

func (s *server) Cancel(ctx context.Context, fr *pb.JobDefinition) (*pb.ReturnMessage, error) {
	return nil, nil
}

func (s *server) RestoreFile(stream pb.Agent_RestoreFileServer) error {
	return nil
}

func (g *GRPC) Client(c *coordinator.Coordinator) (pb.CoordinatorClient, error) {
	if g.clientConn == nil {
		creds := credentials.NewTLS(
			&tls.Config{
				ServerName:   "Coordinator",
				Certificates: []tls.Certificate{g.hostCert.Certificate},
				RootCAs:      g.caCert.Pool,
			})
		var err error
		g.clientConn, err = grpc.Dial(c.Address, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, err
		}
	}

	return pb.NewCoordinatorClient(g.clientConn), nil
}

func SaveConfig(cs config.Store, env map[string]string) error {
	sc := &grpcConfig{}
	hcPath := ""
	hkPath := ""
	caPath := ""
	for k, v := range env {
		switch k {
		case "GRPC_LISTEN":
			sc.listen = v
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
	sc.hostCert, err = certificates.OpenHostCertificate(certificates.CertPath(hcPath), certificates.CertPath(hkPath))
	if err != nil {
		return err
	}

	sc.caCert, err = certificates.OpenCA(certificates.CertPath(caPath))
	if err != nil {
		return err
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
