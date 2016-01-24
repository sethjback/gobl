package apihandler

import "github.com/sethjback/gobble/httpapi"

// Routes are teh agent's routes
var Routes = httpapi.Routes{
	httpapi.Route{
		"Status",
		"GET",
		"/status",
		agentStatus},
	httpapi.Route{
		"Backup",
		"POST",
		"/backups",
		newBackupJob},
	httpapi.Route{
		"Restore",
		"POST",
		"/restores",
		newRestoreJob},
	httpapi.Route{
		"JobList",
		"GET",
		"/jobs",
		jobList},
	httpapi.Route{
		"JobStatus",
		"GET",
		"/jobs/{jobId}",
		jobStatus},
	httpapi.Route{
		"JobCancel",
		"DELETE",
		"/jobs/{jobId}",
		cancelJob},
	httpapi.Route{
		"AgentKey",
		"GET",
		"/key",
		agentKey},
	httpapi.Route{
		"GC",
		"POST",
		"/gc",
		gc}}
