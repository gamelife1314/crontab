package common

import "encoding/json"

// Job represent a executed task
type Job struct {
	// cron job name
	Name string `json:"name"`
	// command
	Command string `json:"command"`
	// cron expr
	CronExpr string `json:"cron_expr"`
}

func (j *Job) String() string {
	if data, err := json.Marshal(j); err != nil {
		return ""
	} else {
		return string(data)
	}
}

// JobLog represent log that will be stored in mongo
type JobLog struct {
	JobName      string `json:"jobName" bson:"jobName"`
	Command      string `json:"command" bson:"command"`
	Err          string `json:"err" bson:"err"`
	PlanTime     string `json:"planTime" bson:"planTime"`
	ScheduleTime int64  `json:"scheduleTime" bson:"scheduleTime"`
	StartTime    int64  `json:"startTime" bson:"startTime"`
	EndTime      int64  `json:"endTime" bson:"endTime"`
}

// JobFilter is used to query log by jobName
type JobFilter struct {
	JobName string `bson:"jobName"`
}

// SortLogByStartTime is used to sort log when query log
type SortLogByStartTime struct {
	SortOrder int `bson:"startTime"`
}

// Response represent http response structure
type Response struct {
	ErrorCode int         `json:"error_code"`
	ErrorMsg  string      `json:"error_msg"`
	Data      interface{} `json:"data"`
}

func (r *Response) String() string {
	if data, err := json.Marshal(r); err != nil {
		return ""
	} else {
		return string(data)
	}
}
