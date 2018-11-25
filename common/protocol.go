package common

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gorhill/cronexpr"
)

// Job represent a executed task
type Job struct {
	// cron job name
	Name string `json:"name"`
	// command
	Command string `json:"command"`
	// cron expr
	CronExpr string `json:"cronExpr"`
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
	OutPut       string `json:"outPut" bson:"outPut"`
	PlanTime     int64  `json:"planTime" bson:"planTime"`
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

// JobSchedulePlan represent job schedule plan
type JobSchedulePlan struct {
	Job      *Job                 `json:"job"`
	Expr     *cronexpr.Expression `json:"expr"`
	NextTime time.Time            `json:"next_time"`
}

func (j *JobSchedulePlan) String() string {
	if data, err := json.Marshal(j); err != nil {
		return ""
	} else {
		return string(data)
	}
}

// JobExecuteInfo
type JobExecuteInfo struct {
	Job        *Job               `json:"job"`
	PlanTime   time.Time          `json:"plan_time"`
	RealTime   time.Time          `json:"real_time"`
	CancelCtx  context.Context    `json:"-"`
	CancelFunc context.CancelFunc `json:"-"`
}

func (j *JobExecuteInfo) String() string {
	if data, err := json.Marshal(j); err != nil {
		return ""
	} else {
		return string(data)
	}
}

// JobEvent
type JobEvent struct {
	EventType int
	Job       *Job
}

func (j *JobEvent) String() string {
	if data, err := json.Marshal(j); err != nil {
		return ""
	} else {
		return string(data)
	}
}

// JobExecuteResult
type JobExecuteResult struct {
	JobExecuteInfo *JobExecuteInfo `json:"job_execute_info"`
	Output         []byte          `json:"output"`
	Err            error           `json:"err"`
	StartTime      time.Time       `json:"start_time"`
	EndTime        time.Time       `json:"end_time"`
}

func (j *JobExecuteResult) String() string {
	if data, err := json.Marshal(j); err != nil {
		return ""
	} else {
		return string(data)
	}
}

// JobExecuteLog
type JobExecuteLog struct {
	JobName      string `json:"jobName" bson:"jobName"`
	Command      string `json:"command" bson:"command"`
	Err          string `json:"err" bson:"err"`
	Output       string `json:"output" bson:"output"`
	PlanTime     int64  `json:"planTime" bson:"planTime"`
	ScheduleTime int64  `json:"scheduleTime" bson:"scheduleTime"`
	StartTime    int64  `json:"startTime" bson:"startTime"`
	EndTime      int64  `json:"endTime" bson:"endTime"`
}

// LogBatch
type LogBatch struct {
	Logs []interface{}
}
