package apihandler

import (
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/gobl/agent/manager"
	"github.com/sethjback/gobl/httpapi"
	"github.com/sethjback/gobl/util/log"
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

func agentKey(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	key, err := manager.Key()
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 500}
	}
	log.Debug("apihsndler", key)

	return httpapi.Response{Data: map[string]interface{}{"keyString": key}, HTTPCode: 200}
}
