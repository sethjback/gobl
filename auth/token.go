package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/sethjback/gobl/goblerr"
)

const (
	ErrorBase64Decode       = "Base64DecodeFailed"
	ErrorJSONEncode         = "JSONEncodeFailed"
	ErrorJWTHeaderDecode    = "JWTHeaderDecodeFailed"
	ErrorJWTHeader          = "JWTHeaderInvalid"
	ErrorJWTClaimsDecode    = "JWTClaimsDecodeFailed"
	ErrorJWTClaims          = "JWTClaimsrInvalid"
	ErrorJWTTokenFormat     = "JWTFormatInvalid"
	ErrorJWTALGInvalid      = "JWTAlgInvalid"
	ErrorJWTSigntureInvalid = "JWTSignatureInvalid"
	ErrorJWTTokenExpired    = "JWTTokenExpired"

	AlgHS256 = "HS256"
	TypeJWT  = "JWT"
)

type Token struct {
	Secret   []byte
	ValidFor int
	Claims   *Claims
	header   *header
}

// Claims holds the claims to be included in the payload of the JWT
type Claims struct {
	Subject    string `json:"sub"`
	Expiration int    `json:"exp"`
	Issued     int    `json:"iat"`
}

type header struct {
	Algorithm string `json:"alg"`
	TokenType string `json:"typ"`
}

func NewToken(secret []byte, validFor int) *Token {
	return &Token{
		Secret:   secret,
		ValidFor: validFor,
		header:   &header{AlgHS256, TypeJWT},
		Claims:   &Claims{}}
}

func (t *Token) Generate() (string, goblerr.Error) {
	h, err := encodePart(t.header)
	if err != nil {
		return "", err
	}

	t.Claims.Issued = int(time.Now().UTC().Unix())
	t.Claims.Expiration = t.Claims.Issued + t.ValidFor
	c, err := encodePart(t.Claims)
	if err != nil {
		return "", err
	}
	combined := append(h, []byte(".")...)
	combined = append(combined, c...)

	signature := generateSignature(combined, t.Secret)

	combined = append(combined, []byte(".")...)
	combined = append(combined, base64Encode(signature)...)

	return string(combined), nil
}

// Parse takes a token string and validates it
func (t *Token) Parse(in string) goblerr.Error {
	split := strings.Split(in, ".")
	if len(split) != 3 {
		return goblerr.New("Token invalid", ErrorJWTTokenFormat, nil)
	}

	sig, err := base64Decode([]byte(split[2]))
	if err != nil {
		return goblerr.New("Token invalid", ErrorJWTTokenFormat, err)
	}

	body := append([]byte(split[0]), []byte(".")...)
	body = append(body, []byte(split[1])...)

	t.header, err = decodeHeader([]byte(split[0]))
	if err != nil {
		return goblerr.New("Token invalid", ErrorJWTTokenFormat, err)
	}

	if t.header.Algorithm != AlgHS256 || t.header.Algorithm != TypeJWT {
		return goblerr.New("Token invalid", ErrorJWTALGInvalid, nil)
	}

	if !verifySignature(body, sig, t.Secret) {
		return goblerr.New("Token invalid", ErrorJWTSigntureInvalid, nil)
	}

	t.Claims, err = decodeClaims([]byte(split[1]))
	if err != nil {
		return goblerr.New("Token invalid", ErrorJWTTokenFormat, err)
	}

	if int(time.Now().Unix()) > t.Claims.Expiration {
		return goblerr.New("Token expired", ErrorJWTTokenExpired, nil)
	}

	return nil
}

func base64Encode(b []byte) []byte {
	buf := make([]byte, base64.RawURLEncoding.EncodedLen(len(b)))
	base64.RawURLEncoding.Encode(buf, b)
	return buf
}

func base64Decode(b []byte) ([]byte, goblerr.Error) {
	buf := make([]byte, base64.RawURLEncoding.DecodedLen(len(b)))
	n, err := base64.RawURLEncoding.Decode(buf, b)
	if err != nil {
		return nil, goblerr.New("Could not decode base64 string", ErrorBase64Decode, err)
	}
	return buf[:n], nil
}

func encodePart(p interface{}) ([]byte, goblerr.Error) {
	JSON, err := json.Marshal(p)
	if err != nil {
		return nil, goblerr.New("Could not json encode token part", ErrorJSONEncode, err)
	}
	return base64Encode(JSON), nil
}

func decodeHeader(raw []byte) (*header, goblerr.Error) {
	s, err := base64Decode(raw)
	if err != nil {
		return nil, goblerr.New("Invalidly encoded header", ErrorJWTHeaderDecode, err)
	}

	var h header
	jerr := json.Unmarshal(s, &h)
	if jerr != nil {
		return nil, goblerr.New("Header not valid", ErrorJWTHeader, jerr)
	}

	return &h, nil
}

func decodeClaims(raw []byte) (*Claims, goblerr.Error) {
	s, err := base64Decode(raw)
	if err != nil {
		return nil, goblerr.New("Invalidly encoded claims", ErrorJWTClaimsDecode, err)
	}

	var c Claims
	jerr := json.Unmarshal(s, &c)
	if jerr != nil {
		return nil, goblerr.New("Claims not valid", ErrorJWTClaims, jerr)
	}

	return &c, nil
}

func generateSignature(body []byte, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(body)
	return mac.Sum(nil)
}

func verifySignature(raw []byte, signature []byte, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(raw)
	return hmac.Equal(signature, mac.Sum(nil))
}
