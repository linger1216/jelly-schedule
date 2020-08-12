package core

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
	"time"
)

type ExecutorConfig struct {
	Interval int `json:"interval" yaml:"interval"`
}

type Executor struct {
	etcd     *Etcd
	db       *sqlx.DB
	ticker   *time.Ticker
	schedule *cron.Cron
}

func NewExecutor(etcd *Etcd, db *sqlx.DB, config ExecutorConfig) *Executor {
	e := &Executor{etcd: etcd, db: db}
	ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
	e.ticker = ticker
	go e.handleTicker()
	e.schedule = cron.New()
	return e
}

func (e *Executor) close() {
	e.schedule.Stop()
}

func (e *Executor) addCronWorkFlow(workflow *WorkFlow) error {
	_, err := e.schedule.AddFunc(workflow.Cron, func() {
		resp, err := e.Exec(workflow)
		if err != nil {
			l.Warnf("%s workflow err:%s", workflow.Name, err.Error())
		}
		l.Debugf("%s workflow response:%v", workflow.Name, resp)
	})
	return err
}

func (e *Executor) handleTicker() {
	for {
		select {
		case <-e.ticker.C:
			l.Debugf("check workflow")
			workFlows, err := e.getAvaiableWorkFLow(1)
			if err != nil {
				panic(err)
			}
			for i := range workFlows {
				err := e.addCronWorkFlow(workFlows[i])
				if err != nil {
					// todo
					// maybe output log
					panic(err)
				}
			}
		}
	}
}

func (e *Executor) getAvaiableWorkFLow(n int) ([]*WorkFlow, error) {
	tx, err := e.db.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	query := getWorkFLowForUpdate(StateAvaiable, 1)
	rows, err := tx.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	workFlows := make([]*WorkFlow, 0)
	for rows.Next() {
		line := make(map[string]interface{})
		err = rows.MapScan(line)
		if err != nil {
			return nil, err
		}
		if tc, err := transWorkflow("", line); tc != nil && err == nil {
			workFlows = append(workFlows, tc)
		}
	}

	if len(workFlows) == 0 {
		return nil, nil
	}

	for i := range workFlows {
		workFlows[i].State = StateExecuting
	}

	query, args, err := upsertWorkflowSql(workFlows)
	if err != nil {
		return nil, err
	}

	l.Debug(query)
	_, err = tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return workFlows, nil
}

func (e *Executor) Exec(workFlow *WorkFlow) (interface{}, error) {
	if workFlow == nil {
		return nil, ErrorInvalidPara
	}
	// 默认的执行方式是串行的
	serialJob := NewSerialJob(nil)
	for _, jobJroup := range workFlow.JobIds {
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
	return serialJob.Exec(context.Background(), workFlow.Para)
}
