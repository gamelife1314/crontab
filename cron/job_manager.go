package cron

import (
	"context"
	"time"

	"github.com/gamelife1314/crontab/common"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

// JobManager
type JobManager struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

// Close
func (j *JobManager) Close() {
	j.client.Close()
}

// watchJobs
func (j *JobManager) watchJobs() (err error) {
	var (
		getResp            *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
		jobEvent           *common.JobEvent
	)

	if getResp, err = j.kv.Get(context.TODO(), common.CronJobDir, clientv3.WithPrefix()); err != nil {
		return
	}
	common.Logger.Infoln("Starting synchronous tasks")
	for _, kvpair = range getResp.Kvs {
		if job, err = common.UnpackJob(kvpair.Value); err == nil {
			jobEvent = common.BuildJobEvent(common.JobSaveEvent, job)
			GlobalScheduler.PushJobEvent(jobEvent)
		}
	}

	go func() {
		watchStartRevision = getResp.Header.Revision + 1
		watchChan = j.watcher.Watch(context.TODO(), common.CronJobDir, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		common.Logger.Infoln("Starting listen etcd for adding and deleting task")
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					jobEvent = common.BuildJobEvent(common.JobSaveEvent, job)
				case mvccpb.DELETE:
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))
					job = &common.Job{Name: jobName}
					jobEvent = common.BuildJobEvent(common.JobDeleteEvent, job)
				}
				GlobalScheduler.PushJobEvent(jobEvent)
			}
		}
	}()
	return
}

// watchKiller
func (j *JobManager) watchKiller() {
	var (
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
		jobEvent   *common.JobEvent
		jobName    string
		job        *common.Job
	)

	go func() {
		watchChan = j.watcher.Watch(context.TODO(), common.CronKillJobDir, clientv3.WithPrefix())
		common.Logger.Infoln("Stating listen tasks those are killed.")
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:
					jobName = common.ExtractKillerName(string(watchEvent.Kv.Key))
					job = &common.Job{Name: jobName}
					jobEvent = common.BuildJobEvent(common.JobKillEvent, job)
					GlobalScheduler.PushJobEvent(jobEvent)
				case mvccpb.DELETE:
				}
			}
		}
	}()
}

var GlobalJobManager *JobManager

func InitJobManager() (err error) {

	var client *clientv3.Client

	config := clientv3.Config{
		Endpoints:   GlobalConfig.Etcd.Endpoints,
		DialTimeout: time.Duration(GlobalConfig.Etcd.DialTimeout) * time.Millisecond,
	}

	if client, err = clientv3.New(config); err != nil {
		return
	}

	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)
	watcher := clientv3.NewWatcher(client)

	GlobalJobManager = &JobManager{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	GlobalJobManager.watchJobs()
	GlobalJobManager.watchKiller()

	return
}

// CreateJobLock
func (j *JobManager) CreateJobLock(jobName string) (jobLock *JobLock) {
	jobLock = InitJobLock(jobName, j.kv, j.lease)
	return
}
