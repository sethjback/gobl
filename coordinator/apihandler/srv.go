package apihandler

import (
	"net/http"
	"runtime"
	"time"

	"github.com/sethjback/gobble/httpapi"
)

func gc(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	runtime.GC()
	return &httpapi.APIResponse{HTTPCode: 200}, nil
}

func coordinatorStatus(w http.ResponseWriter, r *http.Request) (*httpapi.APIResponse, error) {
	status := make(map[string]interface{})

	status["date"] = time.Now().String()

	goR := runtime.NumGoroutine()
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)

	status["goRoutines"] = goR
	status["memory"] = memStat.Alloc

	return &httpapi.APIResponse{Data: status, HTTPCode: 200}, nil
}