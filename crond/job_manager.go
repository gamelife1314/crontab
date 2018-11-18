package crond

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gamelife1314/crontab/common"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

// jobManager is used to describe jobManager struct.
type JobManager struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_JobManager *JobManager
)

// InitJobManager is used to initialize jobManager.
func InitJobManager() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)

	config = clientv3.Config{
		Endpoints:   Config.Etcd.Endpoints,
		DialTimeout: time.Duration(Config.Etcd.DialTimeout) * time.Millisecond,
	}

	if client, err = clientv3.New(config); err != nil {
		return
	}

	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_JobManager = &JobManager{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	return
}

// Close etcd connection
func (j *JobManager) Close() {
	j.client.Close()
}

// SaveJob is used to save job to etcd
func (j *JobManager) SaveJob(job *common.Job) (prevJob *common.Job, err error) {
	var (
		jobKey     string
		jobStrExpr []byte
		putResp    *clientv3.PutResponse
		prevJobObj common.Job
	)

	jobKey = common.CronJobDir + job.Name
	if jobStrExpr, err = json.Marshal(job); err != nil {
		return
	}

	if putResp, err = j.kv.Put(context.TODO(), jobKey, string(jobStrExpr), clientv3.WithPrevKV()); err != nil {
		return
	}

	if putResp.PrevKv != nil {
		if err = json.Unmarshal(putResp.PrevKv.Value, &prevJobObj); err != nil {
			return
		}
		prevJob = &prevJobObj
	}
	return
}

// DeleteJob is used to delete job by key
func (j *JobManager) DeleteJob(name string) (prevJob *common.Job, err error) {
	var (
		jobKey     string
		delResp    *clientv3.DeleteResponse
		prevJobObj common.Job
	)

	jobKey = common.CronJobDir + name
	if delResp, err = j.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); err != nil {
		return
	}

	if len(delResp.PrevKvs) != 0 {
		if err = json.Unmarshal(delResp.PrevKvs[0].Value, &prevJobObj); err != nil {
			return
		}
		prevJob = &prevJobObj
	}
	return
}

// ListJobs is used to list all cron jobs
func (j *JobManager) ListJobs() (jobList []*common.Job, err error) {
	var (
		dirKey  string
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		job     *common.Job
	)
	dirKey = common.CronJobDir

	if getResp, err = j.kv.Get(context.TODO(), dirKey, clientv3.WithPrefix()); err != nil {
		return
	}

	jobList = make([]*common.Job, 0)

	for _, kvPair = range getResp.Kvs {
		job = &common.Job{}
		if err = json.Unmarshal(kvPair.Value, job); err != nil {
			common.Logger.WithFields(logrus.Fields{
				"jobKey":   string(kvPair.Key),
				"jobValue": string(kvPair.Value),
			}).Error("job parse failed!")
			continue
		}
		jobList = append(jobList, job)
	}
	return
}

// KillJob notify worker that execute kill job
func (j *JobManager) KillJob(name string) (err error) {
	var (
		killerKey      string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
	)

	killerKey = common.CronKillJobDir + name

	if leaseGrantResp, err = j.lease.Grant(context.TODO(), 1); err != nil {
		return
	}

	leaseId = leaseGrantResp.ID

	if _, err = j.kv.Put(context.TODO(), killerKey, "", clientv3.WithLease(leaseId)); err != nil {
		return
	}

	return
}
