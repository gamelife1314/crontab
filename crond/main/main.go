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
	"github.com/gamelife1314/crontab/crond"
)

var (
	configFile string
	logFile    string
)

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initArgs() {
	flag.StringVar(&configFile, "config", "./crond.yaml", "set crond config file")
	flag.StringVar(&logFile, "log", "./crond.log", "set crond log file")
	flag.Parse()
}

func init() {
	initEnv()
	initArgs()
	if err := common.SetLogFile(logFile); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if err := crond.LoadConfig(configFile); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if err := crond.InitJobManager(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if err := crond.InitWorkManager(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if err := crond.InitLogManager(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if err := crond.InitApiServer(); err != nil {
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
	crond.G_JobManager.Close()
	crond.G_WorkManager.Close()
	crond.G_LogManager.Close()
	fmt.Println("\r", "Welcome next!")
}

func printCrondInf() {
	var info strings.Builder
	info.WriteString("Crond started at: http://" + crond.Config.Http.Address + ":" + strconv.Itoa(crond.Config.Http.Port))
	info.WriteString("\n")
	info.WriteString("Pid: " + strconv.Itoa(os.Getpid()) + "\n")
	info.WriteString("Please input Ctr+c to terminate crond service.\n")
	info.WriteString("now: " + time.Now().Format("2006/01/02 15:04:05"))
	fmt.Println(info.String())
}
