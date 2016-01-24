package apihandler

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/sethjback/gobble/coordinator/manager"
	"github.com/sethjback/gobble/httpapi"
	"github.com/sethjback/gobble/spec"
)

func restore(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return nil, httpapi.NewError("", "Request too large", 413)
	}

	if err := r.Body.Close(); err != nil {
		return nil, errors.New("")
	}

	restore := new(spec.RestoreRequest)

	err = json.Unmarshal(body, &restore)
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Invalid restore request", 400)
	}

	id, err := manager.RunRestore(restore)

	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Unable to create restore job", 400)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{"id": id}, HTTPCode: 200}, nil
}
