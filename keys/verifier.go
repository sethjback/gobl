package keys

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

type Verifier interface {
	Verify(signed []byte, signature string) error
}

type verifier struct {
	publicKey *rsa.PublicKey
}

func NewVerifier(key *rsa.PublicKey) Verifier {
	return &verifier{publicKey: key}
}

func (v *verifier) Verify(signed []byte, signature string) error {
	dec, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		return errors.New("Could not decode signature string")
	}
	h := sha256.New()
	h.Write(signed)
	d := h.Sum(nil)

	return rsa.VerifyPSS(v.publicKey, crypto.SHA256, d, dec, nil)
}
