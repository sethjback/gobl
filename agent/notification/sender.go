package notification

import (
	"bytes"
	"net/http"

	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/keys"
)

const (
	Retry   = "SenderError"
	Success = "SenderSuccess"
	Fail    = "SenderFail"
)

type Sender struct {
	client  *http.Client
	signer  keys.Signer
	message *Message
}

type Result struct {
	state   string
	err     error
	message *Message
}

func (s *Sender) Do() interface{} {
	req := httpapi.Request{
		Host:   s.message.note.Destination(),
		Path:   "/notification", //TODO: need a place for this
		Client: s.client,
		Method: "POST",
		Body:   bytes.NewReader(s.message.note.Body()),
	}

	sig, err := s.signer.Sign([]byte(req.String()))
	if err != nil {
		return Result{state: Fail, err: err, message: s.message}
	}

	req.AddHeader(httpapi.HeaderGoblSig, sig)

	resp, err := req.Send()
	if err != nil || resp.HTTPCode != 200 {
		return &Result{state: Retry, err: err, message: s.message}
	}

	return &Result{state: Success}
}
