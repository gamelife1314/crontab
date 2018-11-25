package cron

import (
	"github.com/gamelife1314/crontab/common"
	"time"
)

type Scheduler struct {
	jobEventChan      chan *common.JobEvent
	jobPlanTable      map[string]*common.JobSchedulePlan
	jobExecutingTable map[string]*common.JobExecuteInfo
	jobResultChan     chan *common.JobExecuteResult
}

var GlobalScheduler *Scheduler

// PushJobResult
func (s *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	common.Logger.Infoln("Receive new job", jobEvent)
	s.jobEventChan <- jobEvent
}

// PushJobResult
func (s *Scheduler) PushJobResult(result *common.JobExecuteResult) {
	s.jobResultChan <- result
}

func (s *Scheduler) TryStartJob(jobPlan *common.JobSchedulePlan) {
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting   bool
	)
	if jobExecuteInfo, jobExecuting = s.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		return
	}
	jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)
	s.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo
	common.Logger.Infoln("Star task", jobExecuteInfo)
	GlobalExecutor.ExecuteJob(jobExecuteInfo)
}

// TrySchedule
func (s *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan  *common.JobSchedulePlan
		now      time.Time
		nearTime *time.Time
	)
	if len(s.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}
	now = time.Now()

	common.Logger.Infoln("check all tasks, schedule tasks those are expired")
	for _, jobPlan = range s.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			s.TryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now)
		}
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	scheduleAfter = (*nearTime).Sub(now)
	return
}

// handleJobEvent
func (s *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		jobSchedulePlan *common.JobSchedulePlan
		jobExecuteInfo  *common.JobExecuteInfo
		jobExecuting    bool
		jobExisted      bool
		err             error
	)
	switch jobEvent.EventType {
	case common.JobSaveEvent:
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			common.Logger.Errorln("Job saved failed", err.Error())
			return
		}
		s.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case common.JobDeleteEvent:
		common.Logger.Infoln("Delete task：", jobEvent.Job.Name)
		if jobSchedulePlan, jobExisted = s.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(s.jobPlanTable, jobEvent.Job.Name)
		}
	case common.JobKillEvent:
		common.Logger.Infoln("Kill task：", jobEvent.Job.Name)
		if jobExecuteInfo, jobExecuting = s.jobExecutingTable[jobEvent.Job.Name]; jobExecuting {
			jobExecuteInfo.CancelFunc()
		}
	}
}

// handleJobResult
func (s *Scheduler) handleJobResult(result *common.JobExecuteResult) {
	var jobLog *common.JobLog
	delete(s.jobExecutingTable, result.JobExecuteInfo.Job.Name)
	if result.Err != common.LockAlreadyRequiredErr {
		jobLog = &common.JobLog{
			JobName:      result.JobExecuteInfo.Job.Name,
			Command:      result.JobExecuteInfo.Job.Command,
			OutPut:       string(result.Output),
			PlanTime:     result.JobExecuteInfo.PlanTime.UnixNano() / 1000 / 1000,
			ScheduleTime: result.JobExecuteInfo.RealTime.UnixNano() / 1000 / 1000,
			StartTime:    result.JobExecuteInfo.RealTime.UnixNano() / 1000 / 1000,
			EndTime:      result.JobExecuteInfo.RealTime.UnixNano() / 1000 / 1000,
		}
		if result.Err != nil {
			jobLog.Err = result.Err.Error()
		} else {
			jobLog.Err = ""
		}
		GlobalLogSink.Append(jobLog)
	}
}

// scheduleLoop
func (s *Scheduler) scheduleLoop() {
	var (
		jobEvent       *common.JobEvent
		schedulerAfter time.Duration
		scheduleTimer  *time.Timer
		jobResult      *common.JobExecuteResult
	)

	schedulerAfter = s.TrySchedule()
	scheduleTimer = time.NewTimer(schedulerAfter)
	for {
		select {
		case jobEvent = <-s.jobEventChan:
			s.handleJobEvent(jobEvent)
		case <-scheduleTimer.C:
		case jobResult = <-s.jobResultChan:
			s.handleJobResult(jobResult)
		}
		schedulerAfter = s.TrySchedule()
		scheduleTimer.Reset(schedulerAfter)
	}
}

func InitScheduler() (err error) {
	GlobalScheduler = &Scheduler{
		jobEventChan:      make(chan *common.JobEvent, 1000),
		jobPlanTable:      make(map[string]*common.JobSchedulePlan),
		jobExecutingTable: make(map[string]*common.JobExecuteInfo),
		jobResultChan:     make(chan *common.JobExecuteResult, 1000),
	}
	go GlobalScheduler.scheduleLoop()
	return
}
