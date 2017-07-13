package grpcclient

import (
	"crypto/tls"
	"errors"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/sethjback/gobl/certificates"
	"github.com/sethjback/gobl/config"
	pb "github.com/sethjback/gobl/goblgrpc"
)

var coordinatorAddress string
var hostCert *certificates.HostCert
var caCert *certificates.CA
var coordConn *grpc.ClientConn
var connMutex = sync.Mutex{}

func Client() (pb.CoordinatorClient, error) {
	connMutex.Lock()
	defer connMutex.Unlock()
	if coordConn == nil {
		creds := credentials.NewTLS(
			&tls.Config{
				ServerName:   "Coordinator",
				Certificates: []tls.Certificate{hostCert.Certificate},
				RootCAs:      caCert.Pool,
			})
		var err error
		coordConn, err = grpc.Dial(coordinatorAddress, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, err
		}
	}

	return pb.NewCoordinatorClient(coordConn), nil
}

func CloseClient() {
	connMutex.Lock()
	defer connMutex.Unlock()
	if coordConn != nil {
		coordConn.Close()
		coordConn = nil
	}
}

func SaveConfig(cs config.Store, env map[string]string) error {
	hcPath := ""
	hkPath := ""
	caPath := ""
	coordAdd := ""
	for k, v := range env {
		switch k {
		case "HOST_CERT":
			hcPath = v
		case "HOST_KEY":
			hkPath = v
		case "CA_CERT":
			caPath = v
		case "COORDINATOR_ADDRESS":
			coordAdd = v
		}
	}

	if hcPath == "" || hkPath == "" {
		return errors.New("Must provide host certificate and key")
	}
	if caPath == "" {
		return errors.New("Must provide CA certificate")
	}
	if coordAdd == "" {
		return errors.New("Must provide coordinator address")
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

	coordinatorAddress = coordAdd

	return nil
}
