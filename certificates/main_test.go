package certificates

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var caCRT = []byte(`-----BEGIN CERTIFICATE-----
MIIE4jCCAsqgAwIBAgIBATANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDEwZHb2Js
Q0EwHhcNMTcwNjIxMDAwMTA1WhcNMjcwNjIxMDAwMTEyWjARMQ8wDQYDVQQDEwZH
b2JsQ0EwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDtE/yYMGhuqPl/
+ies8P/GOBOo2AJkN2igmpFYoZ3lEL6aLaW6pYp5B+HTwI+uIFZhpMAKREW0n2+K
T1IZu3dKZIPV/DymiSaFtAcUlEN2z2d9dmJBD8Ph9MPjVMcSsP3KfNK2IjgFtqv0
n+Gq/7/M1Dib6INWMTZwTbPrnUaqOQ+X4bs0Xg9GMGgbOIqAsJ7NZTGkkct9Kf/6
QqQ7ozYTxmzHW7VNkxQpT99AwHA0dhqLMsjogIIGsjlaMfLQJklBqRT+oKFguZEG
nSjdl+OYr1m3AXOgVpcHLrHQxNotg4cqvCN/YG3/wkzhNFM5iSi9yne5RWtqKx1e
0Ouj8uNV2UU/0uBJIS/iw0ZOI0ktQYlHVDLEvS+6BBdUMKsHFodW+YEqJf8JjvxX
FDjE+TigitT6Zq+nKI2sqEWm1rkF1kHnl3IZR3TOTc3GSh7Q1y0EsVcnPbyYNX6K
InnNt2OAAGuxUU6gnNw5kciYhq/SCDmgRqmVje93prDC5WTk2HKYVkk1EoXx11i3
yLeYWwEDRH700HBNe+LBPRIa3gh6ZhLff0uSz6yyGHH452t9fe+B59/IjzaLkm+P
q23UVONX7rjHcXFsWXpO6gSUH57jOuRLl45ojSJvxKjqtsTmWkSveqVVlOivT2/8
ToHWT02FYtE7LATG1ourSiQ9QmprGwIDAQABo0UwQzAOBgNVHQ8BAf8EBAMCAQYw
EgYDVR0TAQH/BAgwBgEB/wIBADAdBgNVHQ4EFgQUABTZoVUQUrSjLzAIt2uk8ysh
SdMwDQYJKoZIhvcNAQELBQADggIBAA7h+i9y2DMidt3mO1kFXLCW1ts8IxG6zlTl
tgewYI5+fHMiKlaUItp5A9+KuJEbg27py+KgovsZYKsAHtgdl8VnL5RwYrjbv/tk
aB5hENC/LfFXTzACF7KKQTRe6gG/YWSjySpzhcbmuRo6hCuo+Lv8ekdDpi2Lv7Z+
0OcOZZv19NCsfTNweiTQx7WJEKf6wSbV8LYl9RLMIrXEEtt0hwgFSL4ZAiYvYVPT
xVASvAnC+UWobPPg9w2+u6dcv70BskpoHDLQXXVJHyvi/osfahArKUg6QU++7zG7
Kf8p2w63OFrP6uANWWebVzbdTVCEPBHtxKjOuSaJ9YtrmXEqmoQ9w+bn2NonnxUT
Lwu8K0VPPWKk9OqMJS/Bo314WC/W7BSpxNozubHOCjtsq1+OZaIyR0vOGX+Y6Dub
JbfBu5tdxliYVvqhln4lvxaVAqx9EdYG8dYz59NMsJK+rdgArCUDg5tqBov2s+m5
Eq9/2In9HTbkDQkucW5kDcuWUvhNXVW4lGSriGLiUlcQo0jnjkiuoR4FVfl18lxt
URSaMizSRnxmOG/2OL44cjIfrAUkeDuVAVe/2a4g94eJTZaLthmvObPs7YVfmWdA
vJBvCnqPlnr5wzyRhDi5gjLjVqEgEjY4uKqEKIyuMTU5uZ7LRSAySK0HwR8HSKd1
G9b/A9DC
-----END CERTIFICATE-----
`)

var hostCRT = []byte(`-----BEGIN CERTIFICATE-----
MIIEGzCCAgOgAwIBAgIRAN6GwMOJuE2fGzN/cCt5ivAwDQYJKoZIhvcNAQELBQAw
ETEPMA0GA1UEAxMGR29ibENBMB4XDTE3MDYyMTAwMDI0N1oXDTE5MDYyMTAwMDI0
N1owDjEMMAoGA1UEAxMDMDA3MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAsqfjL5q7mPo691hOmpTdz1g7fRb9ANndoiBSDNw6FkRAh5cmNXKNWg65xEDf
qXxjm4i6c5BYVtTv340y8XSXXcWBPvkVgQWSZ4KGYejwbMINlja9S45+c2thKfqJ
uJAipgMGMeZ1CKqfEH905/iDucRcyX/q5peLO9NEtx7rOIBnXwoLl7RYZqtM8fC4
oXa5PMS++IyH6cnxIJMb5yFh51McXnwxz4OetY/vBZMJRHGWQKMoRjrvZMDSmZd/
x/dF4CpOr8i6p509z3d2f51alngsV5X9cw8HdpzcejwgaTfkZl5+cAQNNw/Kh+tJ
6AvoQgIpTadavLvZyvPHXtoRNwIDAQABo3EwbzAOBgNVHQ8BAf8EBAMCA7gwHQYD
VR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB0GA1UdDgQWBBTsWX5BgA6FFZws
1n3VR+xu4wTsAjAfBgNVHSMEGDAWgBQAFNmhVRBStKMvMAi3a6TzKyFJ0zANBgkq
hkiG9w0BAQsFAAOCAgEACxm+zn94a37Ge41e+RyaYpYLJco77EIdQDC9XMKsQWPq
GJ8Mcx4bUMm+5fTxrIxwL0TugG0H6KuIorkwfZb/uDOKepMqHDMdvVENrLtqcb4y
WP2NVrX6bPPvV+F2qcxxC3Qfhp3+j8Yp54a+h4YHWOH/9tr1PJJfNAnRGu/r2GQf
cybPf5o5vDMwq6blk1EVmKH1IceG2fIAqtDhmhrlnnUVc1r0gxUADbk0PoNQ+x6i
L6O7v+8OdYD2yQN2yAq2AhjwN9BIMFzCflVuZFQn83LQjCddTO7QrlT8TYm389dy
+E+/upakAxkX3B+e4V1BmTCR+VJTtsyp52KSXkP9nIMrucsN18HeO1fiopeMMeW6
hj5+GvO3AuW+ntmctsYvPTZPnjaDxfdOJHMRceCRQ6WAUPuYhvL0dW5Llpw7XEQ/
okXw3dh1Wx3hkGfpbBZk++Q4TioOXsgv7pb7dxTB+69TaQyHPVK249PCoLFfkv9W
g2upu8i2jtsxDyP6SPjkWRpgRmP2UPt3ixsBtcZScv9nUS6smH8nTMlgYI4cGpsN
DDIB2+gFtYiWBfZ1/Y+vhTOHI5jK8MOO0wYdZ3UxJPTobq+5K9I4+t6y8vtWMTMy
RX08DwmPe5miNjI1DisPtZSQ2LGlYP52S7ZeQArIAhDPMG6wfq8UPa/zf4B6aPM=
-----END CERTIFICATE-----
`)

var hostKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAsqfjL5q7mPo691hOmpTdz1g7fRb9ANndoiBSDNw6FkRAh5cm
NXKNWg65xEDfqXxjm4i6c5BYVtTv340y8XSXXcWBPvkVgQWSZ4KGYejwbMINlja9
S45+c2thKfqJuJAipgMGMeZ1CKqfEH905/iDucRcyX/q5peLO9NEtx7rOIBnXwoL
l7RYZqtM8fC4oXa5PMS++IyH6cnxIJMb5yFh51McXnwxz4OetY/vBZMJRHGWQKMo
RjrvZMDSmZd/x/dF4CpOr8i6p509z3d2f51alngsV5X9cw8HdpzcejwgaTfkZl5+
cAQNNw/Kh+tJ6AvoQgIpTadavLvZyvPHXtoRNwIDAQABAoIBAQCmdaVJaquWerhM
VDxQ7ZnKIpSzFaNAkr5d3C13DA8XRhq1+1A/hm9L1OKjiCqdWWfZuEi5emnE5fxm
V8J6lT6fwXGOQjkWESH7TfN18LtrKlfMeU5gwvDxC8DpgyWlEK8n7TNtdSPgolE4
5vj/Vl8tzFcD7CrrFZJGeK/Sy30xE2uH8iKZ6alfNUCHEEUIENKPwuUuvxhUS1gy
rhlF5SkiDFTduhPwEGZRKJ8xNpdKZscwiR1YQbbOY3SzscmJf7X8WvtOEObfOhmi
6VaVomvQyLhOoo+nbwXMXdf1CzjnrrRfIahB9RMqoNudtJNWhcLXr8ouNEXajf2i
NH48GzVZAoGBAOLAR5ytJnHoze+DAKsMysNcj4yP/htgE7n9JqVDVFhAkiefxzQ2
w3Il6/2tCuHLJX5T7kw/18wd5Y0ZR9yVRzgxensAun6IxGa9KGoXKVFQrcgcrxb6
WsT6UiLNb3+b4yK+MtymD8S13izsIjvIsa2Eli1dk22iQXP037KyCrIlAoGBAMmz
arSGyFywNE7SUKQi8gEZwCqiDtp3CkzUZ3y9Y7QoxTp7HybZZm8wPRS54mBk2cTG
jwLVFaDsHYTvh2uktyyeJM8BQUTi4RIDJjIbKqAksqMVypztuCY27zmZEKMnhx+m
qTrvMrLzYtVTxqtlOxdN0EIgeRAtfcN0qds2LAErAoGAW4SForvT33e53mh+VYtF
LxJlsbLQOZZOf+untF33ZeMx2iJH0VAlFCYwGGPGF5nZWSJg9I9z9qM+afOBKItr
gkeeCpUhsD5dHqZL8H3GpFYuvayuElUbW5M4oWlPDi8JvpULjDjN9nP85x6rNnvr
EoCH6GlzPnWVe0qjGsl3Pa0CgYEAmOnpqqVIP2TisyMlOdq4Z/cyxd+IrT10VJzv
PBWFAi+qntR64IQO8Zq7o/vs0LGEm3cBMt+C/zYihwblPsloiW33b+x+pA/xHCvB
CFmqLjDEMXmy2tgqNOaO5LbTcy3jdi5uvBxd7mcwdZSG2KftbZRzn75oqcgjQUwv
/d4K7HMCgYEAme6n26LLhDc11Fa9sBa57lY2K2quQxLrZIELQtaD48kKWK2/c/Vy
bo9pZr5UaQkdKVz5sheFtB4rIJQHs6Ft+/e1Sf+TWjgzdpzReY1Ac2pQBIpWMLsT
o7dcL5PhNOw/NRbfsMjr1apWMj97Ix1fd9WUFvNueTgCCsiI1UJC07k=
-----END RSA PRIVATE KEY-----
`)

func TestNewCA(t *testing.T) {
	assert := assert.New(t)

	ca, err := NewCA(CertPEM(caCRT))
	assert.Nil(err)
	assert.Len(ca.Pool.Subjects(), 1)

	ca, err = NewCA(CertPEM([]byte(`not a cert`)))
	assert.NotNil(err)
	assert.Nil(ca)

	err = ioutil.WriteFile("CAcert", caCRT, 0644)
	if assert.Nil(err) {
		defer os.Remove("CAcert")
	}

	ca, err = NewCA(CertPath("CAcert"))
	assert.Nil(err)
	assert.Len(ca.Pool.Subjects(), 1)

	ca, err = NewCA(CertPath(`not a cert`))
	assert.NotNil(err)
	assert.Nil(ca)
}

func TestNewHostCertificate(t *testing.T) {
	assert := assert.New(t)

	hc, err := NewHostCertificate(CertPEM(hostCRT), CertPEM(hostKey))
	assert.Nil(err)
	assert.NotNil(hc)

	_, err = NewHostCertificate(CertPEM([]byte(`nada cert`)), CertPEM(hostKey))
	assert.NotNil(err)

	_, err = NewHostCertificate(CertPEM(hostCRT), CertPEM([]byte(`nada key`)))
	assert.NotNil(err)

	err = ioutil.WriteFile("hostCRT", hostCRT, 0644)
	if assert.Nil(err) {
		defer os.Remove("hostCRT")
	} else {
		return
	}

	err = ioutil.WriteFile("hostKey", hostKey, 0644)
	if assert.Nil(err) {
		defer os.Remove("hostKey")
	} else {
		return
	}

	hc, err = NewHostCertificate(CertPath("hostCRT"), CertPath("hostKey"))
	assert.Nil(err)
	assert.NotNil(hc)

	hc, err = NewHostCertificate(CertPath("nada"), CertPath("hostKey"))
	assert.Nil(hc)
	assert.NotNil(err)

	hc, err = NewHostCertificate(CertPath("hostCRT"), CertPath("nada"))
	assert.Nil(hc)
	assert.NotNil(err)
}
