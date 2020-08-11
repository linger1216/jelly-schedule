package core

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Executor struct {
	etcd *Etcd
	db   *sqlx.DB
}

func NewExecutor(etcd *Etcd, db *sqlx.DB) *Executor {
	return &Executor{etcd: etcd, db: db}
}

func (e *Executor) Exec(flow *WorkFlow) (interface{}, error) {
	if flow == nil {
		return nil, ErrorInvalidPara
	}
	// 默认的执行方式是串行的
	serialJob := NewSerialJob(nil)
	for _, jobJroup := range flow.JobIds {
		jobs := make([]Job, 0)
		for _, jobId := range jobJroup {
			buf, err := e.etcd.Get(context.Background(), jobId)
			if err != nil {
				return nil, err
			}
			info, err := UnMarshalJobInfo(buf)
			if err != nil {
				return nil, err
			}
			jobs = append(jobs, info.ToJob())
		}
		// 如果某个节点的job数量大于1
		// 说明这个节点可以多个job同时运行
		parallelJob := NewParallelJob(jobs)
		serialJob.Append(parallelJob)
	}
	return serialJob.Exec(context.Background(), flow.Para)
}
