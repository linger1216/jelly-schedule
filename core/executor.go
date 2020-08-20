package core

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/robfig/cron/v3"
	"net/http"
	"strings"
	"time"
)

type ExecutorConfig struct {
	Name                  string `json:"name" yaml:"name" `
	CheckWorkFlowInterval int    `json:"checkWorkFlowInterval" yaml:"checkWorkFlowInterval" `
	MetricPort            int    `json:"metricPort" yaml:"metricPort" `
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
	latency  *Histogram
	success  *Counter
	failed   *Counter

	// stats
	//newWorkFlowStatusCommandQueue
	//statusQueue *WorkFlowStatusCommandQueue
}

type ExecutorContext struct {
	stats *WorkFlowStatus
	entry cron.EntryID
}

// Total duration of requests in seconds.
func initPrometheus(e *Executor, port int) {
	fieldKeys := []string{"name"}

	e.progress = NewGaugeFrom(prometheus.GaugeOpts{
		Namespace: PrometheusNamespace,
		Subsystem: PrometheusSubsystem,
		Name:      "progress",
		Help:      "Workflow or Job progress",
	}, fieldKeys)

	e.latency = NewHistogramFrom(prometheus.HistogramOpts{
		Namespace: PrometheusNamespace,
		Subsystem: PrometheusSubsystem,
		Name:      "latency",
		Help:      "Workflow or Job exec latency",
	}, fieldKeys)

	e.success = NewCounterFrom(prometheus.CounterOpts{
		Namespace: PrometheusNamespace,
		Subsystem: PrometheusSubsystem,
		Name:      "success",
		Help:      "Workflow or Job success count",
	}, fieldKeys)

	e.failed = NewCounterFrom(prometheus.CounterOpts{
		Namespace: PrometheusNamespace,
		Subsystem: PrometheusSubsystem,
		Name:      "failed",
		Help:      "Workflow or Job failed count",
	}, fieldKeys)

	go func() {
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(fmt.Sprintf(":%d", port), m)
	}()
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
	// initPrometheus(e, config.MetricPort)

	// db
	_, err := db.Exec(createWorkflowTableSql())
	if err != nil {
		panic(err)
	}

	go e.handleTicker()

	l.Debugf("exec started.")
	return e
}

func (e *Executor) close() {
	e.workFlowCron.Stop()
}

func (e *Executor) execWorkFlowCron(workflow *WorkFlow) error {
	// 为workflow创建定时任务
	entryId, err := e.workFlowCron.AddFunc(workflow.Cron, func() {
		endpoint := func(ctx context.Context, request interface{}) (response interface{}, err error) {
			return e.exec(workflow)
		}

		etx, ok := e.executorContextMap[workflow.Id]
		if !ok {
			panic("never impossiable")
		}

		now := time.Now()
		ctx := context.Background()
		resp, err := endpoint(ctx, workflow)
		etx.stats.LastExecuteDuration = int64(time.Since(now).Seconds())
		l.Debugf("workflow:%s exec duration:%d", workflow.Name, etx.stats.LastExecuteDuration)
		if err != nil {
			etx.stats.FailedExecuteCount++
			l.Debugf("workflow:%s err:%v", workflow.Name, err.Error())
		} else {
			etx.stats.SuccessExecuteCount++
			l.Debugf("workflow:%s resp:%v", workflow.Name, resp)
		}

		// 成功运行次数
		// -1 代表无限
		if workflow.SuccessLimit > 0 && etx.stats.SuccessExecuteCount >= workflow.SuccessLimit {
			e.workFlowCron.Remove(etx.entry)
			l.Debugf("workflow:%s remove cron:%d", workflow.Name, etx.entry)
			changeWorkFlowState(e.db, StateFinish, workflow)
			return
		}

		// 失败运行次数
		// -1 代表无限
		if workflow.FailedLimit > 0 && etx.stats.FailedExecuteCount >= workflow.FailedLimit {
			e.workFlowCron.Remove(etx.entry)
			l.Debugf("workflow:%s remove cron:%d", workflow.Name, etx.entry)
			changeWorkFlowState(e.db, StateFailed, workflow)
			return
		}
	})

	if err != nil {
		return err
	}

	// 为workflow创建上下文Context
	e.executorContextMap[workflow.Id] = &ExecutorContext{
		stats: &WorkFlowStatus{
			Id: workflow.Id,
		},
		entry: entryId,
	}

	// 运行Cron任务
	e.workFlowCron.Start()
	l.Debugf("workflow:%s add cron:%d", workflow.Name, entryId)
	return nil
}

func (e *Executor) handleTicker() {
	findWorkFlowAndExecCron := func(query string) bool {
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
			err := e.execWorkFlowCron(workFlows[i])
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
			handled = findWorkFlowAndExecCron(getWorkFLowByExecutorBelongForUpdate(e.name, StateAvaiable, 1))
			// 其次查询普通任务
			if !handled {
				handled = findWorkFlowAndExecCron(getWorkFLowForUpdate(StateAvaiable, 1))
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

func (e *Executor) exec(workFlow *WorkFlow) (interface{}, error) {
	if workFlow == nil {
		return nil, ErrorInvalidPara
	}
	finalJob := NewSerialJob(nil)
	// [["a"],["a","b","c"], ["x.y.z"]]
	for _, jobGroup := range workFlow.JobIds {
		parallelJobs := make([]Job, 0)
		for _, group := range jobGroup {
			ids := strings.Split(group, ",")
			if len(ids) > 1 {
				serial := NewSerialJob(nil)
				for _, id := range ids {
					job, err := e.getJob(id)
					if err != nil {
						return nil, err
					}
					serial.Append(job)
				}
				parallelJobs = append(parallelJobs, serial)
			} else if len(ids) == 1 {
				job, err := e.getJob(ids[0])
				if err != nil {
					return nil, err
				}
				parallelJobs = append(parallelJobs, job)
			} else {
				return nil, ErrJobNotFound
			}
		}
		// 如果某个节点的job数量大于1
		// 说明这个节点可以多个job同时运行
		parallelJob := NewParallelJob(parallelJobs)
		finalJob.Append(parallelJob)
	}
	return finalJob.Exec(context.Background(), workFlow.Para)
}

func (e *Executor) getJob(jobId string) (Job, error) {
	buf, err := e.etcd.Get(context.Background(), JobKey(jobId))
	if err != nil {
		return nil, err
	}
	info, err := UnMarshalJobDescription(buf)
	if err != nil {
		return nil, err
	}
	return info.ToJob(), nil
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

/*

func (e *Executor) execWorkFlowCron(workflow *WorkFlow) error {

	// 为workflow创建定时任务
	// entryId, err :=
	e.workFlowCron.AddFunc(workflow.Cron, func() {

		// 要将workflow的执行封装成endpoint
		// 然后剥笋子来封装

		//workflowContext, ok := e.executorContextMap[workflow.Id]
		//if !ok || workflowContext == nil {
		//	panic(fmt.Sprintf("workflowContext %s invalid", workflow.Id))
		//}

		endpoint := func(ctx context.Context, request interface{}) (response interface{}, err error) {
			return e.execByPolicy(nil, workflow)
		}





		// endpoint = Instrumenting(e.latency, e.success, e.failed)(endpoint)

		//
		//state := StateFinish
		//if err != nil {
		//	l.Warnf("%s workflow err:%s", workflow.Name, err.Error())
		//	state = StateFailed
		//}

		//err = changeWorkFlowState(e.db, state, workflow)
		//
		//if err != nil {
		//	//l.Warnf("changeWorkFlowState workflow:%s err:%s", workflow.Name, err.Error())
		//}
		//l.Debugf("%s workflow run success", workflow.Name)

		// 运行次数到了限定
		// 退出
		//if workflowContext.stats.SuccessExecuteCount.Load() >= int32(workflow.SuccessLimit) {
		//	e.workFlowCron.Remove(workflowContext.entry)
		//	return
		//}
	})

	//if err != nil {
	//	return err
	//}

	// 为workflow创建上下文Context
	//e.executorContextMap[workflow.Id] = &ExecutorContext{
	//	stats: &WorkFlowStatus{
	//		Id: workflow.Id,
	//	},
	//	entry: entryId,
	//}

	// 再次运行确认没问题
	// e.workFlowCron.Start()
	return nil
}
*/
