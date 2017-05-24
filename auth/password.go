package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"

	"github.com/sethjback/gobl/goblerr"
)

const (
	ErrorPasswordHash = "PasswordHashFailed"
)

const saltLength = 4

// generateSalt generates a secure random hex encoded salt
func generateSalt() ([]byte, error) {
	b := make([]byte, saltLength/2)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	salt := make([]byte, hex.EncodedLen(len(b)))

	hex.Encode(salt, b)

	return salt, nil
}

// createSaltedHash takes the random salt and the provided password and generates a hex encoded sha1 sum of the combined values
// The returned []byte contans the salt followed by the sha1 sum
func createSaltedHash(salt []byte, password []byte) []byte {
	password = append(salt, password...)

	sha1Sum := sha1.Sum(password)

	hash := make([]byte, hex.EncodedLen(len(sha1Sum)))
	hex.Encode(hash, sha1Sum[:])

	return append(salt, hash[:]...)
}

// PasswordHash generates a hex encoded sha1 hash of the password with crypt/rand salt
// The returned []byte contans the salt followed by the sha1 sum
// The function will only error if the install's secure random generator is not working
func PasswordHash(password []byte) ([]byte, error) {
	salt, err := generateSalt()
	if err != nil {
		return nil, goblerr.New("Unable to generate password hash", ErrorPasswordHash, err)
	}

	return createSaltedHash(salt, password), nil
}

// CheckPassword checks the provided password against the stored password
// It extracts the salt from the front of the stored value, then adds it to the password and re-computes the hash
func CheckPassword(saved []byte, check []byte) bool {
	salt := make([]byte, saltLength)
	copy(salt, saved)

	hash := createSaltedHash(salt, check)

	return bytes.Equal(saved, hash)
}
