package keys

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

// Manager handles storing the instances' key
type Manager struct {
	PrivateKey *rsa.PrivateKey
	PublicKeys map[string]*rsa.PublicKey
}

// OpenPrivateKey reads a pem encode private key from disk
func OpenPrivateKey(path string) (*rsa.PrivateKey, error) {
	privateKeyFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	privateKeyBytes, _ := pem.Decode(privateKeyFile)

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// OpenPublicKey decodes public key file from disk
func OpenPublicKey(path string) (*rsa.PublicKey, error) {
	keyFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	pemBlock, _ := pem.Decode(keyFile)
	return parsePublicKeyBytes(pemBlock.Bytes)
}

// DecodePublicKeyString takes a the base64.URLEncoded DER bytes and converts them to a rsa.PulicKey
func DecodePublicKeyString(keyString string) (*rsa.PublicKey, error) {
	keyBytes, err := base64.URLEncoding.DecodeString(keyString)
	if err != nil {
		return nil, err
	}
	return parsePublicKeyBytes(keyBytes)
}

func parsePublicKeyBytes(b []byte) (*rsa.PublicKey, error) {
	pI, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}

	pKey, ok := pI.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("Invalid public key")
	}

	return pKey, nil
}

// PublicKey exports the public key (DER) as a base64.URLEncoded string
func (m *Manager) PublicKey() (string, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(&m.PrivateKey.PublicKey)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(pubASN1), nil
}

// Sign uses the private key to sign provided string and returns the result as a base64.URLEncoded string
func (m *Manager) Sign(s string) (string, error) {
	h := sha256.New()
	h.Write([]byte(s))
	d := h.Sum(nil)
	sig, err := rsa.SignPSS(rand.Reader, m.PrivateKey, crypto.SHA256, d, nil)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(sig), nil
}

// KeyForHost gets a publick key for the given host
func (m *Manager) KeyForHost(host string) (*rsa.PublicKey, error) {
	key, ok := m.PublicKeys[host]
	if !ok {
		return nil, errors.New("Could not find the key for that host")
	}

	return key, nil
}

// VerifySignature validates the signature string
func VerifySignature(key *rsa.PublicKey, signed []byte, signature string) error {
	dec, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		return errors.New("Could not decode signature string")
	}

	h := sha256.New()
	h.Write(signed)
	d := h.Sum(nil)

	err = rsa.VerifyPSS(key, crypto.SHA256, d, dec, nil)

	if err != nil {
		return err
	}

	return nil
}
