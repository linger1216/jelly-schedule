package core

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	jsoniter "github.com/json-iterator/go"
	"github.com/linger1216/jelly-schedule/utils"
	"github.com/linger1216/jelly-schedule/utils/snowflake"
	"github.com/scylladb/termtables"
	"github.com/valyala/fasttemplate"
	"os"
	"time"
)

var (
	JobPrefix = `/schedule/job`
	JobFormat = fasttemplate.New(JobPrefix+`/{Id}`, "{", "}")
	TTL       = int64(10)
)

type JobStats struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	ServicePath string `json:"servicePath"`
	JobPath     string `json:"jobPath"`
}

func (w JobStats) String() string {
	table := termtables.CreateTable()
	table.AddHeaders("Field", "Value")
	table.AddRow("Name", w.Name)
	table.AddRow("Host", w.Host)
	table.AddRow("Port", w.Port)
	table.AddRow("ServicePath", w.ServicePath)
	table.AddRow("JobPath", w.JobPath)
	return table.Render()
}

func JobKey(id string) string {
	s := JobFormat.ExecuteString(map[string]interface{}{
		"Id": id,
	})
	return s
}

func MarshalJobStats(j *JobStats) ([]byte, error) {
	return jsoniter.ConfigFastest.Marshal(j)
}

func UnMarshalJobStats(buf []byte) (*JobStats, error) {
	s := &JobStats{}
	err := jsoniter.ConfigFastest.Unmarshal(buf, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

type JobServer struct {
	stats   JobStats
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

	l.Debugf("job %s started", ret.stats.Name)
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
	jsonBuf, _ := MarshalJobStats(&w.stats)
	return w.etcd.InsertKV(context.Background(), JobKey(w.stats.Id), string(jsonBuf), w.leaseId)
}

func (w *JobServer) Stats() string {
	return w.stats.String()
}

func (w *JobServer) Close() {
	w.ticker.Stop()
}
