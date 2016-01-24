package sqlite

const (
	agentsTable      = "agents"
	jobsTable        = "jobs"
	filesTable       = "files"
	definitionsTable = "definitions"
	schedulesTable   = "schedules"

	createAgentsTable      = "CREATE TABLE IF NOT EXISTS `agents` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `name` string, `address` string, `publicKey` string)"
	createDefinitionsTable = "CREATE TABLE IF NOT EXISTS `definitions` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `agent` int, `paramiters` blob)"
	createJobsTable        = "CREATE TABLE IF NOT EXISTS `jobs` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `jobtype` string, `agent` int, `definition` blob, `start` datetime, `end` datetime, `state` int, `message` string)"
	createFilesTable       = "CREATE TABLE IF NOT EXISTS `files` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `job` int, `state` int, `message` string, `path` string, `name` string, `signature` blob, `meta` blob)"
	createSchedulesTable   = "CREATE TABLE IF NOT EXISTS `schedules` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `backup` int, `schedule` blob)"
)
