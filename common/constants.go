package common

const (
	CronJobDir     = "/cron/job/"
	CronKillJobDir = "/cron/killed/job/"
	CronWorkerDir  = "/cron/workers/"
	CronLockDir    = "/cron/lock/"
)

const (
	_ = iota
	JobSaveEvent
	JobDeleteEvent
	JobKillEvent
)
