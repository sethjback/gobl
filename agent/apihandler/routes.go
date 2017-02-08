package apihandler

import "github.com/sethjback/gobl/httpapi"

var Routes = []httpapi.Route{

	// Agent
	httpapi.Route{
		Method:  "GET",
		Path:    "/status",
		Handler: agentStatus,
	},
	httpapi.Route{
		Method:  "GET",
		Path:    "/key",
		Handler: agentKey,
	},

	// Jobs
	httpapi.Route{
		Method:  "POST",
		Path:    "/backups",
		Handler: newBackupJob,
	},
	httpapi.Route{
		Method:  "POST",
		Path:    "/restores",
		Handler: newRestoreJob,
	},
	httpapi.Route{
		Method:  "GET",
		Path:    "/jobs",
		Handler: jobList,
	},
	httpapi.Route{
		Method:  "GET",
		Path:    "/jobs/:id",
		Handler: jobStatus,
	},
	httpapi.Route{
		Method:  "DELETE",
		Path:    "/jobs/:id",
		Handler: cancelJob,
	},
}
