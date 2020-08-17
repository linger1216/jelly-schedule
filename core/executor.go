package core

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"time"
)

type ExecutorConfig struct {
	Name                  string `json:"name" yaml:"name" `
	CheckWorkFlowInterval int    `json:"checkWorkFlowInterval" yaml:"checkWorkFlowInterval" `
}

type Executor struct {
	name string
	etcd *Etcd
	db   *sqlx.DB
	// 用于检查是否有可用的workflow
	CheckWorkFlowTicker *time.Ticker
	// 用于执行workflow具体的cron
	workFlowCron       *cron.Cron
	executorContextMap map[string]*ExecutorContext

	// prometheus
	// 进度 gauge
	// 响应时间
	// 成功次数
	// 失败次数
	progress *Gauge
}

type ExecutorContext struct {
	stats *WorkFlowStats
	entry cron.EntryID
}

func NewExecutor(etcd *Etcd, db *sqlx.DB, config ExecutorConfig) *Executor {

	e := &Executor{etcd: etcd, db: db}
	ticker := time.NewTicker(time.Duration(config.CheckWorkFlowInterval) * time.Second)
	e.CheckWorkFlowTicker = ticker
	e.workFlowCron = cron.New()

	// todo
	// 需要确认在这里合适与否
	e.workFlowCron.Start()

	e.name = config.Name
	e.executorContextMap = make(map[string]*ExecutorContext)

	// prometheus
	fieldKeys := []string{"name", "access_key", "error"}
	e.progress = NewGaugeFrom(prometheus.GaugeOpts{}, fieldKeys)

	go e.handleTicker()
	return e
}

func (e *Executor) close() {
	e.workFlowCron.Stop()
}

func (e *Executor) addCronWorkFlow(workflow *WorkFlow) error {

	// 为workflow创建定时任务
	entryId, err := e.workFlowCron.AddFunc(workflow.Cron, func() {

		workflowContext, ok := e.executorContextMap[workflow.Id]
		if !ok || workflowContext == nil {
			panic(fmt.Sprintf("workflowContext %s invalid", workflow.Id))
		}

		_, err := e.execByPolicy(workflowContext.stats, workflow)

		state := StateFinish
		if err != nil {
			l.Warnf("%s workflow err:%s", workflow.Name, err.Error())
			state = StateFailed
		}

		err = changeWorkFlowState(e.db, state, workflow)

		if err != nil {
			//l.Warnf("changeWorkFlowState workflow:%s err:%s", workflow.Name, err.Error())
		}
		//l.Debugf("%s workflow run success", workflow.Name)

		// 运行次数到了限定
		// 退出
		if workflowContext.stats.SuccessExecuteCount.Load() >= int32(workflow.ExecuteLimit) {
			e.workFlowCron.Remove(workflowContext.entry)
			return
		}
	})
	if err != nil {
		return err
	}

	// 为workflow创建上下文Context
	e.executorContextMap[workflow.Id] = &ExecutorContext{
		stats: &WorkFlowStats{
			Id: workflow.Id,
		},
		entry: entryId,
	}

	// 再次运行确认没问题
	// e.workFlowCron.Start()
	return nil
}

func (e *Executor) handleTicker() {

	addCronWorkFlow := func(query string) bool {
		if len(query) == 0 {
			return false
		}
		workFlows, err := e.getAvaiableWorkFLow(query)
		if err != nil {
			panic(err)
		}
		if len(workFlows) == 0 {
			return false
		}
		for i := range workFlows {
			err := e.addCronWorkFlow(workFlows[i])
			if err != nil {
				panic(err)
			}
		}
		return true
	}

	for {
		select {
		case <-e.CheckWorkFlowTicker.C:
			// 首先查询自己专属的任务
			var handled bool
			var query string
			query = getWorkFLowByExecutorBelongForUpdate(e.name, StateAvaiable, 1)
			handled = addCronWorkFlow(query)

			// 其次查询普通任务
			if !handled {
				query = getWorkFLowForUpdate(StateAvaiable, 1)
				handled = addCronWorkFlow(query)
			}
		}
	}
}

func (e *Executor) getAvaiableWorkFLow(query string) ([]*WorkFlow, error) {
	tx, err := e.db.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback()
	}()

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

	query, args, err := upsertWorkflowSql(workFlows...)
	if err != nil {
		return nil, err
	}

	//l.Debug(query)
	_, err = tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return workFlows, nil
}

func (e *Executor) execByPolicy(stats *WorkFlowStats, workFlow *WorkFlow) (interface{}, error) {
	retry := 1
	if workFlow.ErrorPolicy == ErrPolicyRetry {
		retry = DefaultRetryCount
	}

	var resp interface{}
	var err error

	for retry > 0 {
		t := time.Now()
		resp, err = e.exec(workFlow)
		stats.LastExecuteDuration.CAS(stats.LastExecuteDuration.Load(), int32(time.Since(t).Milliseconds()))
		if err != nil {
			stats.FailedExecuteCount.Inc()
			continue
		}
		stats.SuccessExecuteCount.Inc()
		retry--
	}

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (e *Executor) exec(workFlow *WorkFlow) (interface{}, error) {
	if workFlow == nil {
		return nil, ErrorInvalidPara
	}

	// 默认的执行方式是串行的
	serialJob := NewSerialJob(nil)
	for _, jobJroup := range workFlow.JobIds {
		jobs := make([]Job, 0)
		for _, jobId := range jobJroup {
			buf, err := e.etcd.Get(context.Background(), JobKey(jobId))
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

func changeWorkFlowState(db *sqlx.DB, state string, workflow *WorkFlow) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback()
	}()

	workflow.State = state
	query, args, err := upsertWorkflowSql(workflow)
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
