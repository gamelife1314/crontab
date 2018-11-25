package cron

import (
	"context"
	"time"

	"github.com/gamelife1314/crontab/common"
	"go.etcd.io/etcd/clientv3"
)

type Register struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease

	localIP string
}

var GlobalRegister *Register

// Close
func (r *Register) Close() {
	r.client.Close()
}

// keepOnline
func (r *Register) keepOnline() {
	var (
		regKey         string
		leaseGrantResp *clientv3.LeaseGrantResponse
		err            error
		keepAliveChan  <-chan *clientv3.LeaseKeepAliveResponse
		keepAliveResp  *clientv3.LeaseKeepAliveResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
	)

	for {
		common.Logger.Infoln("Try to register client", r.localIP)
		regKey = common.CronWorkerDir + r.localIP
		cancelFunc = nil
		if leaseGrantResp, err = r.lease.Grant(context.TODO(), 10); err != nil {
			goto RETRY
		}
		if keepAliveChan, err = r.lease.KeepAlive(context.TODO(), leaseGrantResp.ID); err != nil {
			goto RETRY
		}
		cancelCtx, cancelFunc = context.WithCancel(context.TODO())
		if _, err = r.kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseGrantResp.ID)); err != nil {
			goto RETRY
		}
		for {
			select {
			case keepAliveResp = <-keepAliveChan:
				if keepAliveResp == nil {
					goto RETRY
				}
			}
		}

	RETRY:
		time.Sleep(1 * time.Second)
		if cancelFunc != nil {
			cancelFunc()
		}
	}
}

// InitRegister
func InitRegister() (err error) {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		kv      clientv3.KV
		lease   clientv3.Lease
		localIp string
	)
	config = clientv3.Config{
		Endpoints:   GlobalConfig.Etcd.Endpoints,
		DialTimeout: time.Duration(GlobalConfig.Etcd.DialTimeout) * time.Millisecond,
	}
	if client, err = clientv3.New(config); err != nil {
		return
	}
	if localIp, err = common.GetLocalIP(); err != nil {
		return
	}
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	GlobalRegister = &Register{
		client:  client,
		kv:      kv,
		lease:   lease,
		localIP: localIp,
	}
	go GlobalRegister.keepOnline()
	return
}
