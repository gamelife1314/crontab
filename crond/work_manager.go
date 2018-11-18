package crond

import (
	"context"
	"time"

	"github.com/gamelife1314/crontab/common"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

type WorkManager struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var G_WorkManager *WorkManager

// InitWorkManager is used to init G_WorkManager
func InitWorkManager() (err error) {
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

	G_WorkManager = &WorkManager{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

// Close etcd connection
func (w *WorkManager) Close() {
	w.client.Close()
}

// ListWorkers is used to list all workers
func (w *WorkManager) ListWorkers() (workers []string, err error) {
	var (
		getResp  *clientv3.GetResponse
		kv       *mvccpb.KeyValue
		workerIp string
	)
	workers = make([]string, 0)
	if getResp, err = w.kv.Get(context.TODO(), common.CronWorkerDir, clientv3.WithPrefix()); err != nil {
		return
	}
	for _, kv = range getResp.Kvs {
		workerIp = common.ExtractWorkerIp(string(kv.Key))
		workers = append(workers, workerIp)
	}
	return
}
