package apihandler

import (
	"runtime"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sethjback/gobl/coordinator/email"
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
	err := email.SendEmail("This is a test email from gobl. Let me be the first to congratulate you on receiving this message: it means your email is configured correctly. Way to go!", "Gobl Coordinator")
	if err != nil {
		return httpapi.Response{Error: err, HTTPCode: 400}
	}
	return httpapi.Response{HTTPCode: 200}
}
