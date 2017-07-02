package apihandler

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/gobl/certificates"
	"github.com/sethjback/gobl/coordinator/manager"
	"github.com/sethjback/gobl/gobldb/errors"
	"github.com/sethjback/gobl/goblerr"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/model"
)

func newCAKey(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	k := &model.Key{}
	err := r.JsonBody(k)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	if k.Certificate == "" && k.Key == "" {
		k = nil
	} else {
		_, err = certificates.OpenHostCertificate(certificates.CertPEM([]byte(k.Certificate)), certificates.CertPEM([]byte(k.Key)))
		if err != nil {
			return httpapi.Response{Error: err, HTTPCode: 400}
		}
	}

	err = manager.NewCAKey(k)
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}

	return httpapi.Response{HTTPCode: 201}
}

func getCertificate(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	k, err := manager.GetCAKey()

	if err != nil {
		gerr := err.(*goblerr.Error)
		HTTPCode := 400
		if gerr.Code == errors.ErrCodeNotFound {
			HTTPCode = 200
		}
		return httpapi.Response{Error: err, HTTPCode: HTTPCode}
	}

	return httpapi.Response{HTTPCode: 200, Data: map[string]interface{}{"certificate": k.Certificate}}
}
