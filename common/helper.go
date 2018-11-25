package common

import (
	"context"
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"net"
	"strings"
	"time"
)

// ExtractWorkerIp extract ip from worker key
func ExtractWorkerIp(rawIp string) string {
	return strings.TrimPrefix(rawIp, CronWorkerDir)
}

// ExtractKillerName
func ExtractKillerName(killerKey string) string {
	return strings.TrimPrefix(killerKey, CronKillJobDir)
}

// ExtractJobName
func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, CronJobDir)
}

// BuildResponse build http response
func BuildResponse(errorNum int, msg string, data interface{}) (resp []byte, err error) {
	var (
		response Response
	)

	response.ErrorCode = errorNum
	response.ErrorMsg = msg
	response.Data = data

	resp, err = json.Marshal(response)

	return
}

// UnpackJob
func UnpackJob(value []byte) (job *Job, err error) {
	job = &Job{}
	err = json.Unmarshal(value, job)
	return
}

// BuildJobEvent
func BuildJobEvent(eventType int, job *Job) *JobEvent {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

// BuildJobSchedulePlan
func BuildJobSchedulePlan(job *Job) (jobSchedulePlan *JobSchedulePlan, err error) {
	var cronExpr *cronexpr.Expression
	if cronExpr, err = cronexpr.Parse(job.CronExpr); err != nil {
		return
	}

	jobSchedulePlan = &JobSchedulePlan{
		Job:      job,
		Expr:     cronExpr,
		NextTime: cronExpr.Next(time.Now()),
	}
	return
}

// BuildJobExecuteInfo
func BuildJobExecuteInfo(jobSchedulePlan *JobSchedulePlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobSchedulePlan.Job,
		PlanTime: jobSchedulePlan.NextTime,
		RealTime: time.Now(),
	}
	jobExecuteInfo.CancelCtx, jobExecuteInfo.CancelFunc = context.WithCancel(context.TODO())
	return
}

// GetLocalIP
func GetLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet
		isIpNet bool
	)

	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}

	for _, addr = range addrs {
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				return
			}
		}
	}
	err = NetworkLocalIpNotFound
	return
}
