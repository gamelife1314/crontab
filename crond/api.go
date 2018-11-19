package crond

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gamelife1314/crontab/common"
)

type ApiServer struct {
	httpServer *http.Server
}

var (
	Server *ApiServer
)

// InitApiServer initialize api server
func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux
		listener   net.Listener
		address    string
		httpServer *http.Server
	)

	// config routes
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", HandleJobSave)
	mux.HandleFunc("/job/delete", HandleJobDelete)
	mux.HandleFunc("/job/list", HandleJobList)
	mux.HandleFunc("/job/kill", HandleJobKill)
	mux.HandleFunc("/job/log", HandleJobLog)
	mux.HandleFunc("/worker/list", HandleWorkerList)

	address = Config.Http.Address + ":" + strconv.Itoa(Config.Http.Port)
	if listener, err = net.Listen("tcp", address); err != nil {
		return err
	}

	httpServer = &http.Server{
		ReadTimeout:  time.Duration(Config.Http.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(Config.Http.WriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	Server = &ApiServer{
		httpServer: httpServer,
	}

	go httpServer.Serve(listener)

	return
}

func handleError(response http.ResponseWriter) {
	if err := recover(); err != nil {
		errStr := fmt.Errorf("%v", err)
		if bytes, err := common.BuildResponse(-1, errStr.Error(), nil); err == nil {
			response.Write(bytes)
		} else {
			common.Logger.Fatal("A fatal error occur.")
		}

	}
}

// handleJobSave is used to handle requests of creating or updating.
func HandleJobSave(response http.ResponseWriter, request *http.Request) {
	defer handleError(response)
	var (
		err                        error
		jobName, command, cronExpr string
		job, prevJob               *common.Job
		bytes                      []byte
	)

	if err = request.ParseForm(); err != nil {
		panic(err.Error())
	}

	jobName = request.PostForm.Get("jobName")
	command = request.PostForm.Get("command")
	cronExpr = request.PostForm.Get("cronExpr")

	job = &common.Job{
		Name:     jobName,
		Command:  command,
		CronExpr: cronExpr,
	}

	if prevJob, err = G_JobManager.SaveJob(job); err != nil {
		panic(err.Error())
	}

	if bytes, err = common.BuildResponse(0, "success", prevJob); err == nil {
		response.Write(bytes)
	}
	return
}

// handleJobDelete is used to handle requests of delete
func HandleJobDelete(response http.ResponseWriter, request *http.Request) {
	defer handleError(response)
	var (
		err     error
		name    string
		prevJob *common.Job
		bytes   []byte
	)

	if err = request.ParseForm(); err != nil {
		panic(err.Error())
	}

	name = request.PostForm.Get("jobName")

	if prevJob, err = G_JobManager.DeleteJob(name); err != nil {
		panic(err.Error())
	}

	if bytes, err = common.BuildResponse(0, "success", prevJob); err == nil {
		response.Write(bytes)
	}
	return
}

// handleJobList is used to list all jobs.
func HandleJobList(response http.ResponseWriter, request *http.Request) {
	defer handleError(response)
	var (
		jobList []*common.Job
		bytes   []byte
		err     error
	)

	if jobList, err = G_JobManager.ListJobs(); err != nil {
		panic(err.Error())
	}

	if bytes, err = common.BuildResponse(0, "success", jobList); err == nil {
		response.Write(bytes)
	}

	return
}

// handleJobKill is used to kill someone job.
func HandleJobKill(response http.ResponseWriter, request *http.Request) {
	defer handleError(response)
	var (
		err   error
		name  string
		bytes []byte
	)

	if err = request.ParseForm(); err != nil {
		panic(err.Error())
	}

	name = request.PostForm.Get("jobName")

	if err = G_JobManager.KillJob(name); err != nil {
		panic(err.Error())
	}

	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
		response.Write(bytes)
	}

	return
}

// handleJobLog is used to query task logs.
func HandleJobLog(response http.ResponseWriter, request *http.Request) {
	defer handleError(response)
	var (
		err                   error
		name                  string
		limit, skip           int
		limitParam, skipParam string
		logArr                []*common.JobLog
		bytes                 []byte
	)

	if err = request.ParseForm(); err != nil {
		panic(err.Error())
	}
	name = request.PostForm.Get("jobName")
	if name == "" {
		panic("jobName is required!")
	}
	skipParam = request.PostForm.Get("skip")
	if skipParam == "" {
		skipParam = "0"
	}
	limitParam = request.PostForm.Get("limit")
	if limitParam == "" {
		limitParam = "10"
	}

	if skip, err = strconv.Atoi(skipParam); err != nil {
		panic(err.Error())
	}

	if limit, err = strconv.Atoi(limitParam); err != nil {
		panic(err.Error())
	}

	if logArr, err = G_LogManager.ListLog(name, int64(skip), int64(limit)); err != nil {
		panic(err.Error())
	}

	if bytes, err = common.BuildResponse(0, "success", logArr); err == nil {
		response.Write(bytes)
	}

	return
}

// handleWorkerList is used to list worker nodes.
func HandleWorkerList(response http.ResponseWriter, request *http.Request) {
	defer handleError(response)
	var (
		workers []string
		err     error
		bytes   []byte
	)

	if workers, err = G_WorkManager.ListWorkers(); err != nil {
		panic(err.Error())
	}

	if bytes, err = common.BuildResponse(0, "success", workers); err == nil {
		response.Write(bytes)
	}

	return
}
