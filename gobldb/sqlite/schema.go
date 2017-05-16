package sqlite

const (
	agentsTable      = "agents"
	jobsTable        = "jobs"
	filesTable       = "files"
	definitionsTable = "definitions"
	schedulesTable   = "schedules"
	usersTable       = "users"

	createAgentsTable      = "CREATE TABLE IF NOT EXISTS `agents` (`_id` INTEGER PRIMARY KEY AUTOINCREMENT, `id` string, `name` string, `address` string, `publicKey` string)"
	createDefinitionsTable = "CREATE TABLE IF NOT EXISTS `definitions` (`_id` INTEGER PRIMARY KEY AUTOINCREMENT, `id` string, `data` blob)"
	createJobsTable        = "CREATE TABLE IF NOT EXISTS `jobs` (`_id` INTEGER PRIMARY KEY AUTOINCREMENT, `id` string, `agent` integer, `start` datetime, `end` datetime, `state`, string, `data` blob)"
	createFilesTable       = "CREATE TABLE IF NOT EXISTS `files` (`_id` INTEGER PRIMARY KEY AUTOINCREMENT, `job` string, `state` string, `error` blob, `file` blob, `level` int, `parent` string, `name` string)"
	createSchedulesTable   = "CREATE TABLE IF NOT EXISTS `schedules` (`_id` INTEGER PRIMARY KEY AUTOINCREMENT, `id` string, `schedule` blob)"
	createUsersTable       = "CREATE TABLE IF NOT EXISTS `users` (`_id` INTEGER PRIMARY KEY AUTOINCREMENT, `email` string, `password` string, `lastlogin` datetime, CHECK(email <> ''))"
)
