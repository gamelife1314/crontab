package cron

import (
	"context"
	"github.com/gamelife1314/crontab/common"
	"go.etcd.io/etcd/clientv3"
)

// JobLock
type JobLock struct {
	kv         clientv3.KV
	lease      clientv3.Lease
	jobName    string
	cancelFunc context.CancelFunc
	leaseId    clientv3.LeaseID
	isLocked   bool
}

// InitJobLock
func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock) {
	jobLock = &JobLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
	}
	return
}

// TryLock
func (j *JobLock) TryLock() (err error) {
	var leaseGrantResponse *clientv3.LeaseGrantResponse
	var cancelCtx context.Context
	var cancelFunc context.CancelFunc
	var leaseId clientv3.LeaseID
	var leaseKeepAliveChan <-chan *clientv3.LeaseKeepAliveResponse

	if leaseGrantResponse, err = j.lease.Grant(context.TODO(), 5); err != nil {
		return
	}
	leaseId = leaseGrantResponse.ID
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())
	if leaseKeepAliveChan, err = j.lease.KeepAlive(cancelCtx, leaseId); err != nil {
		cancelFunc()
		j.lease.Revoke(context.TODO(), leaseId)
		return
	}

	go func() {
		var leaseKeepResp *clientv3.LeaseKeepAliveResponse
		for {
			select {
			case leaseKeepResp = <-leaseKeepAliveChan:
				if leaseKeepResp != nil {
					goto END
				}
			}
		}
	END:
	}()

	var txn clientv3.Txn
	var lockKey string
	var txnResp *clientv3.TxnResponse
	txn = j.kv.Txn(context.TODO())
	lockKey = common.CronLockDir + j.jobName
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))
	if txnResp, err = txn.Commit(); err != nil {
		cancelFunc()
		j.lease.Revoke(context.TODO(), leaseId)
		return
	}
	if !txnResp.Succeeded {
		cancelFunc()
		j.lease.Revoke(context.TODO(), leaseId)
		err = common.LockAlreadyRequiredErr
		return
	}

	j.leaseId = leaseId
	j.cancelFunc = cancelFunc
	j.isLocked = true
	return
}

// Unlock
func (j *JobLock) Unlock() {
	if j.isLocked {
		j.cancelFunc()
		j.lease.Revoke(context.TODO(), j.leaseId)
	}
}
