package manager

import (
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/sethjback/gobl/config"
	"github.com/sethjback/gobl/keys"
	"github.com/sethjback/gobl/model"
	"github.com/sethjback/gobl/util/log"
	"github.com/stretchr/testify/assert"
)

var testPrivateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAvh7k7d7LufSmk4AuVqPljDbnsf58iYHfuN9HLq2FHr+EsY55
2kQgCoX0SIqBf8CiYxk77zxxuolz5qgPwZB5u3Yi/lLPfaVoZLo96WrNKjukFubJ
NBvlbPTS4/YNPD90BRKYmbKURDZ2xvTPnizTAGy5Jk0BHZn4kaef+43gAUNuaxCO
YuKzoxRc31kE1l0Wls5pcSni9XIWhAgTn5X2jfopk+bGgcs9Tu7ZZJ/uTfSdL/cV
iMYL9I5nFYncWpPqcBKY4C60C8KGYDVdj/gkfmhALx1Wmqb6S5Iyc75EFdr9ASqX
XCWFz9+OyjkiTtArWnGZbjWvUd9voA4vqBI95wIDAQABAoIBAD4cFanoGSIc3LZf
L8Q6MumFnle1zbWWaiIZP0XuhgivhIgFBaXaj6Ugcdeo9/lmUyaQvdXAJ19LPEPk
L5GKw1oMlA4Fu6dOfDY76bHxpCjh5w9cQer2GhNoP+UdIuHF0P8/Pf8oKevG5zLE
E3eXKS+AVVQ/39dtz5i17Dvf84g1kRbe+Zj2J4tqCLL/SXmvGQctsY3Kbtup7qsz
5a09IVAYNwRoAfT0Y6fKrxUosu8ReG/46g6gCxzEoESiQjVNOwag5+fHIAnp+EkQ
45VKI4GFoEYDXJmK7DsaPNi5nfE9DSozLVY24GzqJrm2zuDabbhIAazfJoGGasom
fAULd8ECgYEA4OohwW1G4sY4JuswJjzOZvarQvmp+jxo9bL7UOk06RVKPlUpvKb0
+hReC8lh7sDrWqNhb92nGhqF1xW/t2gCyUpbhQ1DUbHR0HzM6rDjB6MxHtNEU97m
+5apV5+AAflPBB7zkDU2lZYZQIstqL/21KO/AGHTE79D/ySZf0C6eDcCgYEA2GWx
pZLe10bgUpcodRRTRmZ+7HUbv83/s7sjRvGOh9M1IX+nIP38iq6LdblgYrOOkBXP
qYxLYVRPJPf8wqRf6RIPAMe3IlhKRCohuiIXzjT/xY6ZAmNGD1OfINyhqTf8/Dd2
6G4G/P+PmDbM0BtmcuA9WQUa/hcs+gFV3KBcL9ECgYEAse3HQoQ2ndR+O5u3faiB
CMd/eP6Vz9bWmfk8BChZqUMkdudcm1fhWa2fMOfhx8Vq60txG9RYC7iLxTn5bxij
i6Z9fGafqRNpjuwMGGZTVIlvpJkx5r/iL4pi8WTHGcinD/WEbcMLKY+S7pKsTmF+
3X2k7qJ1H5wiKMhFfnwwiEUCgYEAlexW2KzZwPJ05iOdvwfW7haC5xX3pLp25rHH
rhYbNpUo4U2Mn/n35qkpK+XEFn3qTn8eAYyWiRcdQjKhpsS1Qkflpxe3FI9w3KsH
9Oo77fygG+JAtfvLhUDdJapWQmPs3V0b/8qDAvOYK9ADtEbXgs1DE5LK8bFi0s7s
Cs/7LpECgYASgKrUuulgsGBJ4vNR9DnrUV5Fplzkq6H3Ji2BH0cNghzjw1BlnOPJ
DR2YkKGSx4Nuy+asbHaehmNVXcgckKX5a18eqP8XxVQAjnGn4nFiWoLBMrHlTT+z
l4c9e94mec2bPMhzrDILgtWCTzoShttvG/jR3sSZb3mxYySIHs6Ahw==
-----END RSA PRIVATE KEY-----`)

func TestGetAgentKey(t *testing.T) {
	assert := assert.New(t)

	log.Init(config.Log{Level: log.Level.Warn})

	pkb, _ := pem.Decode(testPrivateKey)
	pk, err := x509.ParsePKCS1PrivateKey(pkb.Bytes)
	if !assert.Nil(err) {
		return
	}

	pks, err := keys.PublicKey(pk)
	if !assert.Nil(err) {
		return
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ps := strings.Split(r.URL.Path, "/")[1:]
		switch ps[0] {
		case "s":
			if err != nil {
				w.WriteHeader(400)
				w.Write([]byte(`{"message":"` + err.Error() + `"}`))
			} else {
				w.WriteHeader(200)
				w.Write([]byte(`{"data":{"keyString":"` + pks + `"}}`))
			}
		case "f":
			w.WriteHeader(400)
			w.Write([]byte(`{"message":"error"}`))
		}
	}))
	defer ts.Close()

	a := model.Agent{
		ID:      uuid.New().String(),
		Address: ts.URL + "/s",
	}

	k, err := getAgentKey(a, keys.NewSigner(pk))
	assert.Nil(err)
	assert.NotEmpty(k)
}
