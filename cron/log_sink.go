package cron

import (
	"context"
	"github.com/gamelife1314/crontab/common"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"
)

// LogSink
type LogSink struct {
	client         *mongo.Client
	logCollection  *mongo.Collection
	logChan        chan *common.JobLog
	autoCommitChan chan *common.LogBatch
}

var GlobalLogSink *LogSink

func InitLogSink() (err error) {
	var client *mongo.Client

	if client, err = mongo.Connect(context.TODO(), GlobalConfig.Mongo.Url); err != nil {
		return
	}

	GlobalLogSink = &LogSink{
		client:         client,
		logCollection:  client.Database(GlobalConfig.Mongo.Database).Collection(GlobalConfig.Mongo.Collection),
		logChan:        make(chan *common.JobLog, 1000),
		autoCommitChan: make(chan *common.LogBatch, 1000),
	}

	go GlobalLogSink.writeLoop()

	return
}

// Close
func (l *LogSink) Close() {
	l.client.Disconnect(context.TODO())
}

// Append
func (l *LogSink) Append(jobLog *common.JobLog) {
	select {
	case l.logChan <- jobLog:
	default:
		// drop log
	}
}

func (l *LogSink) saveLogs(batch *common.LogBatch) {
	l.logCollection.InsertMany(context.TODO(), batch.Logs)
	common.Logger.Infoln("Save logs to mongo")
}

//  writeLoop
func (l *LogSink) writeLoop() {
	var log *common.JobLog
	var logBatch *common.LogBatch
	var commitTimer *time.Timer
	var timeoutBatch *common.LogBatch
	for {
		select {
		case log = <-l.logChan:
			if logBatch == nil {
				logBatch = &common.LogBatch{}
				commitTimer = time.AfterFunc(time.Duration(GlobalConfig.Client.JobLogCommitTimeout)*time.Millisecond, func(batch *common.LogBatch) func() {
					return func() {
						l.autoCommitChan <- batch
					}
				}(logBatch))
			}
			logBatch.Logs = append(logBatch.Logs, log)
			if len(logBatch.Logs) >= GlobalConfig.Client.JobLogBatchSize {
				l.saveLogs(logBatch)
				logBatch = nil
				commitTimer.Stop()
			}
		case timeoutBatch = <-l.autoCommitChan:
			if timeoutBatch != logBatch || len(timeoutBatch.Logs) == 0 {
				continue
			}
			l.saveLogs(timeoutBatch)
			logBatch = nil
		}
	}
}
