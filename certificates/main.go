package certificates

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"

	"github.com/sethjback/gobl/model"
	"github.com/square/certstrap/pkix"
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

func OpenCA(opt CertOption) (ca *CA, err error) {
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

func OpenHostCertificate(cert, key CertOption) (hostCert *HostCert, err error) {
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

func NewCACertificate() (*model.Key, error) {
	key, err := pkix.CreateRSAKey(4086)
	if err != nil {
		return nil, err
	}

	cert, err := pkix.CreateCertificateAuthority(
		key,
		"",
		20,
		"Gobl",
		"USA",
		"",
		"",
		"GobleCA")

	if err != nil {
		return nil, err
	}

	keystring, err := key.ExportPrivate()
	if err != nil {
		return nil, err
	}

	certstring, err := cert.Export()
	if err != nil {
		return nil, err
	}

	return &model.Key{Key: string(keystring), Certificate: string(certstring)}, nil
}

func NewHostCertificate(CA model.Key, name string) (*model.Key, error) {
	caCrt, err := pkix.NewCertificateFromPEM([]byte(CA.Certificate))
	if err != nil {
		return nil, err
	}

	caKey, err := pkix.NewKeyFromPrivateKeyPEM([]byte(CA.Key))
	if err != nil {
		return nil, err
	}

	key, err := pkix.CreateRSAKey(4086)
	if err != nil {
		return nil, err
	}

	csr, err := pkix.CreateCertificateSigningRequest(key, "", nil, nil, "Gobl", "", "", "", name)
	if err != nil {
		return nil, err
	}

	cert, err := pkix.CreateCertificateHost(caCrt, caKey, csr, 20)
	if err != nil {
		return nil, err
	}

	keybytes, err := key.ExportPrivate()
	if err != nil {
		return nil, err
	}

	certbytes, err := cert.Export()
	if err != nil {
		return nil, err
	}

	return &model.Key{Key: string(keybytes), Certificate: string(certbytes)}, nil
}
