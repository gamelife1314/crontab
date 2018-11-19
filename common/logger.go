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

// SetLogFile set log file path.
func SetLogFile(filename string) (error) {
	var (
		file *os.File
		err  error
		flag int
	)
	flag = os.O_WRONLY|os.O_APPEND
	if filename == "" {
		filename = "crond.log"
		flag = flag | os.O_CREATE
	}
	if file, err = os.OpenFile(filename, flag, 0666); err != nil {
		return err
	}
	Logger.Out = file

	return nil
}
