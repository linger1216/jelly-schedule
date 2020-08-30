package core

import (
	"context"
	"github.com/linger1216/jelly-schedule/utils"
	"github.com/linger1216/jelly-schedule/utils/snowflake"
	"github.com/valyala/fasttemplate"
	"os"
)

// 提供给用户使用, 内部会调用RPC, 抽象成服务
// 并注册到etcd
var (
	JobPrefix = `/schedule/job`
	JobFormat = fasttemplate.New(JobPrefix+`/{Name}`, "{", "}")
	TTL       = int64(10)
)

func genJobKey(id string) string {
	s := JobFormat.ExecuteString(map[string]interface{}{
		"Name": id,
	})
	return s
}

type JobServer struct {
	stats JobDescription
	job   Job
	etcd  *Etcd
}

func NewJobServer(etcd *Etcd, id string, job Job) *JobServer {
	ret := &JobServer{}
	jobPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if len(id) == 0 {
		id = snowflake.Generate()
	}

	ret.stats.Id = id
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

	serve := newJsonRPCServer(ret.stats, ret.job)
	go func() {
		err := serve.Start()
		if err != nil {
			panic(err)
		}
	}()

	l.Debugf("job %s started: host:%s port:%d id:%s", ret.stats.Name, ret.stats.Host, port, ret.stats.Id)
	return ret
}

func (w *JobServer) register() error {
	jsonBuf, _ := MarshalJobDescription(&w.stats)
	return w.etcd.KeepaliveWithTTL(context.Background(), genJobKey(w.stats.Id), string(jsonBuf), TTL)
}

func (w *JobServer) Stats() string {
	return w.stats.String()
}

func (w *JobServer) Close() {
}
