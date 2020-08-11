package core

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/linger1216/jelly-schedule/utils"
	"github.com/linger1216/jelly-schedule/utils/snowflake"
	"github.com/valyala/fasttemplate"
	"os"
	"time"
)

// 提供给用户使用, 内部会调用RPC, 抽象成服务
// 并注册到etcd
var (
	JobPrefix = `/schedule/job`
	JobFormat = fasttemplate.New(JobPrefix+`/{Id}`, "{", "}")
	TTL       = int64(10)
)

func JobKey(id string) string {
	s := JobFormat.ExecuteString(map[string]interface{}{
		"Id": id,
	})
	return s
}

type JobServer struct {
	stats   JobInfo
	job     Job
	etcd    *Etcd
	leaseId clientv3.LeaseID
	ticker  *time.Ticker
}

func NewJobServer(etcd *Etcd, job Job) *JobServer {
	ret := &JobServer{}
	jobPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	ret.stats.Id = snowflake.Generate()
	ret.stats.JobPath = jobPath
	ret.stats.Name = job.Name()
	ret.stats.Host = utils.GetHost()
	ret.stats.ServicePath = `rpc`
	port, err := utils.GetFreePort()
	if err != nil {
		panic(err)
	}
	ret.stats.Port = port

	ret.job = job
	ret.etcd = etcd

	err = ret.register()
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Duration(TTL/2) * time.Second)
	ret.ticker = ticker
	go ret.handleTicker()

	serve := newJsonRPCServer(ret.stats, ret.job)
	go func() {
		err := serve.Start()
		if err != nil {
			panic(err)
		}
	}()

	l.Debugf("job %s started: %d", ret.stats.Name, port)
	return ret
}

func (w *JobServer) handleTicker() {
	for {
		select {
		case <-w.ticker.C:
			if err := w.etcd.RenewLease(context.Background(), w.leaseId); err != nil {
				l.Debugf("renew %s lease err:%s", w.stats.Name, err.Error())
			} else {
				//l.Debugf("renew %s lease ok", w.stats.Name)
			}
		}
	}
}

func (w *JobServer) register() error {
	if w.leaseId == 0 {
		roleLeaseId, err := w.etcd.GrantLease(TTL)
		if err != nil {
			return err
		}
		w.leaseId = roleLeaseId
	}
	jsonBuf, _ := MarshalJobInfo(&w.stats)
	return w.etcd.InsertKV(context.Background(), JobKey(w.stats.Id), string(jsonBuf), w.leaseId)
}

func (w *JobServer) Stats() string {
	return w.stats.String()
}

func (w *JobServer) Close() {
	w.ticker.Stop()
}
