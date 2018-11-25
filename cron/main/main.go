package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gamelife1314/crontab/common"
	"github.com/gamelife1314/crontab/cron"
)

var (
	configFile string
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initArgs() {
	flag.StringVar(&configFile, "config", "./cron.yaml", "set crond config file")
	flag.Parse()
}

func init() {
	initEnv()
	initArgs()
	var err error
	common.Logger.SetOutput(os.Stdout)
	if err = cron.InitConfig(configFile); err != nil {
		goto ERROR
	}
	if err = cron.InitRegister(); err != nil {
		goto ERROR
	}
	if err = cron.InitLogSink(); err != nil {
		goto ERROR
	}
	if err = cron.InitExecutor(); err != nil {
		goto ERROR
	}
	if err = cron.InitScheduler(); err != nil {
		goto ERROR
	}
	if err = cron.InitJobManager(); err != nil {
		goto ERROR
	}
ERROR:
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func main() {
	var (
		signalChan chan os.Signal
	)
	signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	printCrondInf()
	<-signalChan
	cron.GlobalJobManager.Close()
	cron.GlobalLogSink.Close()
	cron.GlobalJobManager.Close()
	fmt.Println("\r", "Welcome next!")
}

func printCrondInf() {
	var info strings.Builder
	info.WriteString("Cron started \n")
	info.WriteString("Pid: " + strconv.Itoa(os.Getpid()) + "\n")
	info.WriteString("Please input Ctr+c to terminate cron service.\n")
	info.WriteString("now: " + time.Now().Format("2006/01/02 15:04:05"))
	fmt.Println(info.String())
}
