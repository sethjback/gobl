package keys

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
)

type Signer interface {
	Sign(input []byte) (string, error)
}

type signer struct {
	privateKey *rsa.PrivateKey
}

func NewSigner(key *rsa.PrivateKey) Signer {
	return &signer{privateKey: key}
}

func (s *signer) Sign(input []byte) (string, error) {
	h := sha256.New()
	h.Write(input)
	d := h.Sum(nil)
	sig, err := rsa.SignPSS(rand.Reader, s.privateKey, crypto.SHA256, d, nil)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(sig), nil
}
