package keys

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

type Signer interface {
	Sign(input []byte) (string, error)
}

type Verifier interface {
	Verify(signed []byte, signature string) error
}

type signer struct {
	privateKey *rsa.PrivateKey
}

type verifier struct {
	publicKey *rsa.PublicKey
}

func NewSigner(key *rsa.PrivateKey) Signer {
	return &signer{privateKey: key}
}

func NewVerifier(key *rsa.PublicKey) Verifier {
	return &verifier{publicKey: key}
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

func (v *verifier) Verify(signed []byte, signature string) error {
	dec, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		return errors.New("Could not decode signature string")
	}
	h := sha256.New()
	h.Write(signed)
	d := h.Sum(nil)

	if err = rsa.VerifyPSS(v.publicKey, crypto.SHA256, d, dec, nil); err != nil {
		return err
	}

	return nil
}
