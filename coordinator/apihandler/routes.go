package apihandler

import "github.com/sethjback/gobl/httpapi"

// Routes are the agent's routes
var Routes = []httpapi.Route{
	httpapi.Route{
		Method:  "GET",
		Path:    "/status",
		Handler: coordinatorStatus},

	//
	//AGENTS
	//

	httpapi.Route{
		Method:  "GET",
		Path:    "/agents",
		Handler: agentList},

	httpapi.Route{
		Method:  "POST",
		Path:    "/agents",
		Handler: addAgent},

	httpapi.Route{
		Method:  "GET",
		Path:    "/agents/:id",
		Handler: getAgent},

	httpapi.Route{
		Method:  "GET",
		Path:    "/agents/:id/status",
		Handler: agentStatus},

	httpapi.Route{
		Method:  "PUT",
		Path:    "/agents/:id",
		Handler: updateAgent},

	//
	//JOBS
	//

	httpapi.Route{
		Method:  "GET",
		Path:    "/jobs/:id",
		Handler: jobStatus},

	httpapi.Route{
		Method:  "POST",
		Path:    "/jobs/:id/files",
		Handler: addJobFile},

	httpapi.Route{
		Method:  "POST",
		Path:    "/jobs/:id/complete",
		Handler: finishJob},

	httpapi.Route{
		Method:  "GET",
		Path:    "/jobs",
		Handler: jobList},

	httpapi.Route{
		Method:  "GET",
		Path:    "/jobs/:id/files",
		Handler: jobFiles},

	httpapi.Route{
		Method:  "GET",
		Path:    "/jobs/:id/directories",
		Handler: jobDirectories},

	//
	// JOB DEFINITIONS
	//
	httpapi.Route{
		Method:  "GET",
		Path:    "/job-definitions",
		Handler: jobDefinitionList},

	httpapi.Route{
		Method:  "GET",
		Path:    "/job-definitions/:id",
		Handler: getJobDefinition},

	httpapi.Route{
		Method:  "DELETE",
		Path:    "/job-definitions/:id",
		Handler: deleteJobDefinition},

	httpapi.Route{
		Method:  "PUT",
		Path:    "/job-definitions/:id",
		Handler: updateJobDefinition},

	httpapi.Route{
		Method:  "POST",
		Path:    "/job-definitions",
		Handler: createJobDefinition},

	//
	// SCHEDULES
	//
	httpapi.Route{
		Method:  "GET",
		Path:    "/schedules",
		Handler: scheduleList},

	httpapi.Route{
		Method:  "POST",
		Path:    "/schedules",
		Handler: addSchedule},

	httpapi.Route{
		Method:  "DELETE",
		Path:    "/schedules/{sID}",
		Handler: deleteSchedule},

	httpapi.Route{
		Method:  "PUT",
		Path:    "/schedules/{sID}",
		Handler: updateSchedule},

	httpapi.Route{
		Method:  "POST",
		Path:    "/email",
		Handler: testEmail},
}
