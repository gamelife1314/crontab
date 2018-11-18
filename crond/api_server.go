package crond

import (
	"net"
	"net/http"
	"strconv"
	"time"
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
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)
	mux.HandleFunc("/job/log", handleJobLog)
	mux.HandleFunc("/worker/list", handleWorkerList)

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

// handleJobSave is used to handle requests of creating or updating.
func handleJobSave(response http.ResponseWriter, request *http.Request) {

}

// handleJobDelete is used to handle requests of delete
func handleJobDelete(response http.ResponseWriter, request *http.Request) {

}

// handleJobList is used to list all jobs.
func handleJobList(response http.ResponseWriter, request *http.Request) {

}

// handleJobKill is used to kill someone job.
func handleJobKill(response http.ResponseWriter, request *http.Request) {

}

// handleJobLog is used to query task logs.
func handleJobLog(response http.ResponseWriter, request *http.Request) {

}

// handleWorkerList is used to list worker nodes.
func handleWorkerList(response http.ResponseWriter, request *http.Request) {

}
