package common

type Job struct {
	Name     string
	Command  string
	CronExpr string
}

type JobLog struct {
	JobName      string `json:"jobName" bson:"jobName"`
	Command      string `json:"command" bson:"command"`
	Err          string `json:"err" bson:"err"`
	PlanTime     string `json:"planTime" bson:"planTime"`
	ScheduleTime int64  `json:"scheduleTime" bson:"scheduleTime"`
	StartTime    int64  `json:"startTime" bson:"startTime"`
	EndTime      int64  `json:"endTime" bson:"endTime"`
}

type JobFilter struct {
	JobName string `bson:"jobName"`
}

type SortLogByStartTime struct {
	SortOrder int `bson:"startTime"`
}

type Response struct {
	ErrorCode int         `json:"error_code"`
	ErrorMsg  string      `json:"error_msg"`
	Data      interface{} `json:"data"`
}
