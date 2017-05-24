package notification

import (
	"bytes"
	"fmt"
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
		Headers: http.Header{},
		Host:    s.message.note.Host(),
		Path:    s.message.note.Path(),
		Client:  s.client,
		Method:  "POST",
		Body:    bytes.NewReader(s.message.note.Body()),
	}
	fmt.Printf("Sending: %+v\n", req)

	resp, err := req.Send(s.signer)
	if err != nil || resp.HTTPCode != 200 {
		fmt.Printf("err: %+v\n", resp, err)
		return &Result{state: Retry, err: err, message: s.message}
	}

	return &Result{state: Success}
}
