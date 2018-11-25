package cron

import (
	"github.com/gamelife1314/crontab/common"
	"math/rand"
	"os/exec"
	"time"
)

type Executor struct {
}

var GlobalExecutor *Executor

func InitExecutor() (err error) {
	GlobalExecutor = &Executor{}
	return
}

func (e *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	common.Logger.Infoln("Starting to execute task")
	go func() {
		var result *common.JobExecuteResult
		var jobLock *JobLock
		var err error
		result = &common.JobExecuteResult{
			JobExecuteInfo: info,
			Output:         make([]byte, 0),
		}
		jobLock = GlobalJobManager.CreateJobLock(info.Job.Name)
		result.StartTime = time.Now()
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		err = jobLock.TryLock()
		defer jobLock.Unlock()
		if err != nil {
			result.Err = err
			result.EndTime = time.Now()
		} else {
			var cmd *exec.Cmd
			var output []byte
			result.StartTime = time.Now()
			cmd = exec.CommandContext(info.CancelCtx, GlobalConfig.Client.BashExecutePath, "-c", info.Job.Command)
			output, err = cmd.CombinedOutput()
			result.EndTime = time.Now()
			result.Output = output
			result.Err = err
		}
		GlobalScheduler.PushJobResult(result)
	}()
}
