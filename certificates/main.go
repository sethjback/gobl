package certificates

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
)

type CertOption interface{}

type CertPath string

type CertPEM []byte

type CA struct {
	Pool *x509.CertPool
}

type HostCert struct {
	Certificate tls.Certificate
}

func NewCA(opt CertOption) (ca *CA, err error) {
	var pem []byte
	if o, ok := opt.(CertPath); ok {
		pem, err = ioutil.ReadFile(string(o))
		if err != nil {
			return nil, err
		}
	} else if o, ok := opt.(CertPEM); ok {
		pem = []byte(o)
	}

	ca = &CA{Pool: x509.NewCertPool()}
	if ok := ca.Pool.AppendCertsFromPEM(pem); !ok {
		return nil, errors.New("Unable to add pem to certificat pool")
	}
	return ca, nil
}

func NewHostCertificate(cert, key CertOption) (hostCert *HostCert, err error) {
	var cPem []byte
	var kPem []byte
	if o, ok := cert.(CertPath); ok {
		cPem, err = ioutil.ReadFile(string(o))
		if err != nil {
			return nil, err
		}
	} else if o, ok := cert.(CertPEM); ok {
		cPem = []byte(o)
	}

	if o, ok := key.(CertPath); ok {
		kPem, err = ioutil.ReadFile(string(o))
		if err != nil {
			return nil, err
		}
	} else if o, ok := key.(CertPEM); ok {
		kPem = []byte(o)
	}

	hostCert = &HostCert{}
	hostCert.Certificate, err = tls.X509KeyPair(cPem, kPem)
	return hostCert, err
}
