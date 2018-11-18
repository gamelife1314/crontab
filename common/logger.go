package common

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger = logrus.New()

func init() {
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func SetLogFile(filename string) {
	var (
		file *os.File
		err  error
	)
	if filename == "" {
		filename = "crond.log"
	}
	if file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		Logger.Error(err.Error())
		os.Exit(-1)
	}
	Logger.Out = file
}
