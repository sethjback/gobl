package apihandler

import (
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/gobl/agent/manager"
	"github.com/sethjback/gobl/httpapi"
)

func agentStatus(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	status := make(map[string]interface{})
	status["date"] = time.Now().String()
	mstate := manager.Status()
	for k, v := range mstate {
		status[k] = v
	}

	return httpapi.Response{Data: status, HTTPCode: 200}
}
