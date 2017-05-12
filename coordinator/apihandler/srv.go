package apihandler

import (
	"runtime"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/gobl/coordinator/manager"
	"github.com/sethjback/gobl/httpapi"
)

func coordinatorStatus(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	status := make(map[string]interface{})

	status["date"] = time.Now().String()

	goR := runtime.NumGoroutine()
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)

	status["goRoutines"] = goR
	status["memory"] = memStat.Alloc

	return httpapi.Response{Data: status, HTTPCode: 200}
}

func testEmail(r *httpapi.Request, ps httprouter.Params) httpapi.Response {
	err := manager.SendTestEmail()
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}
	return httpapi.Response{HTTPCode: 200}
}
