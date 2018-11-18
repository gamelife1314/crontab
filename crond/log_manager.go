package crond

import (
	"context"
	"time"

	"github.com/gamelife1314/crontab/common"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/options"
)

type LogManager struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

var G_LogManager *LogManager

func InitLogManager() (err error) {
	var (
		client *mongo.Client
	)

	if client, err = mongo.Connect(context.TODO(), Config.Mongo.Url); err != nil {
		return
	}

	G_LogManager = &LogManager{
		client:        client,
		logCollection: client.Database(Config.Mongo.Database).Collection(Config.Mongo.Collection),
	}

	return
}

func (l *LogManager) Close() {
	var (
		ctx context.Context
	)
	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	l.client.Disconnect(ctx)
}

// ListLog is used to read log from mongo
func (l *LogManager) ListLog(name string, skip, limit int64) (logArr []*common.JobLog, err error) {

	var (
		filter  *common.JobFilter
		logSort *common.SortLogByStartTime
		cursor  mongo.Cursor
		jobLog  *common.JobLog
		findOpt *options.FindOptions
	)

	logArr = make([]*common.JobLog, 0)
	filter = &common.JobFilter{
		JobName: name,
	}

	logSort = &common.SortLogByStartTime{SortOrder: -1}

	findOpt = &options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort:  logSort,
	}

	if cursor, err = l.logCollection.Find(context.TODO(), filter, findOpt); err != nil {
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		jobLog = &common.JobLog{}
		if err = cursor.Decode(jobLog); err != nil {
			continue
		}
		logArr = append(logArr, jobLog)
	}

	return
}
