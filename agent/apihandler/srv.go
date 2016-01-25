package apihandler

import (
	"net/http"
	"runtime"
	"time"

	"github.com/sethjback/gobl/agent/manager"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/util/log"
)

func gc(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	runtime.GC()
	return &httpapi.APIResponse{HTTPCode: 200}, nil
}

func agentStatus(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	log.Debug("apihandler", "agentStatus called")
	status := make(map[string]interface{})

	status["date"] = time.Now().String()

	status["status"] = manager.Status()

	return &httpapi.APIResponse{Data: status, HTTPCode: 200}, nil
}

func agentKey(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	key, err := manager.GetKey()
	if err != nil {
		return nil, httpapi.NewError(err.Error(), "Server Error", 500)
	}

	return &httpapi.APIResponse{Data: map[string]interface{}{"keyString": key}, HTTPCode: 200}, nil

}
