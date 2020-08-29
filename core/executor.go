package core

import (
	"context"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/jmoiron/sqlx"
	"github.com/linger1216/jelly-schedule/parser"
	"github.com/robfig/cron/v3"
	"time"
)

type ExecutorConfig struct {
	Name                  string `json:"name" yaml:"name" `
	CheckWorkFlowInterval int    `json:"checkWorkFlowInterval" yaml:"checkWorkFlowInterval" `
	MetricPort            int    `json:"metricPort" yaml:"metricPort" `
	Separate              string `json:"separate,omitempty" yaml:"separate" `
}

type Executor struct {
	name                string
	etcd                *Etcd
	db                  *sqlx.DB
	CheckWorkFlowTicker *time.Ticker
	workFlowCron        *cron.Cron
	contexts            *SyncMap
	exprListener        *ExprListener
	separate            string
}

type ExecutorContext struct {
	stats *WorkFlowStatus
	entry cron.EntryID
}

func NewExecutor(etcd *Etcd, db *sqlx.DB, config ExecutorConfig) *Executor {

	e := &Executor{etcd: etcd, db: db}
	ticker := time.NewTicker(time.Duration(config.CheckWorkFlowInterval) * time.Second)
	e.CheckWorkFlowTicker = ticker
	e.workFlowCron = cron.New()
	e.name = config.Name
	e.separate = config.Separate
	if len(e.separate) == 0 {
		e.separate = ";"
	}
	e.contexts = NewSyncMap()
	e.exprListener = NewExprListener(e.getJob, e.andJob, e.orJob)
	_, err := db.Exec(createWorkflowTableSql())
	if err != nil {
		panic(err)
	}

	e.workFlowCron.Start()
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
		var ctx *ExecutorContext
		if v, ok := e.contexts.Get(workflow.Id); ok {
			ctx = v.(*ExecutorContext)
		}

		if ctx == nil {
			panic("never impossiable")
		}

		if ctx.stats.Executing {
			l.Debugf("workflow:%s executing, ignore...", workflow.Name)
			ctx.stats.MaxExecuteCount++
			if ctx.stats.MaxExecuteCount >= 64 {
				e.workFlowCron.Remove(ctx.entry)
				l.Debugf("workflow:%s remove cron:%d", workflow.Name, ctx.entry)
				err := changeWorkFlowState(e.db, StateFailed, workflow)
				if err != nil {
					l.Debugf("workflow:%s changeWorkFlowState err:%s", workflow.Name, err.Error())
				}
			}
			return
		}

		ctx.stats.Executing = true
		defer func() {
			ctx.stats.Executing = false
		}()

		now := time.Now()
		resp, err := e.exec(workflow)
		ctx.stats.LastExecuteDuration = int64(time.Since(now).Seconds())
		l.Debugf("workflow:%s exec duration:%d", workflow.Name, ctx.stats.LastExecuteDuration)
		if err != nil {
			ctx.stats.FailedExecuteCount++
			l.Debugf("workflow:%s err:%v", workflow.Name, err.Error())
		} else {
			ctx.stats.SuccessExecuteCount++
			l.Debugf("workflow:%s resp:%v", workflow.Name, resp)
		}

		// 成功运行次数
		// -1 代表无限
		if workflow.SuccessLimit > 0 && ctx.stats.SuccessExecuteCount >= workflow.SuccessLimit {
			e.workFlowCron.Remove(ctx.entry)
			l.Debugf("workflow:%s remove cron:%d", workflow.Name, ctx.entry)
			err = changeWorkFlowState(e.db, StateFinish, workflow)
			if err != nil {
				l.Debugf("workflow:%s changeWorkFlowState err:%s", workflow.Name, err.Error())
			}
			return
		}

		// 失败运行次数
		// -1 代表无限
		if workflow.FailedLimit > 0 && ctx.stats.FailedExecuteCount >= workflow.FailedLimit {
			e.workFlowCron.Remove(ctx.entry)
			l.Debugf("workflow:%s remove cron:%d", workflow.Name, ctx.entry)
			err = changeWorkFlowState(e.db, StateFailed, workflow)
			if err != nil {
				l.Debugf("workflow:%s changeWorkFlowState err:%s", workflow.Name, err.Error())
			}
			return
		}
	})
	if err != nil {
		return err
	}

	// 为workflow创建上下文Context
	e.contexts.Put(workflow.Id, &ExecutorContext{
		stats: &WorkFlowStatus{
			Id:        workflow.Id,
			Executing: false,
		},
		entry: entryId,
	})

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

func (e *Executor) exec(workFlow *WorkFlow) (string, error) {
	if workFlow == nil {
		return "", ErrorInvalidPara
	}
	is := antlr.NewInputStream(workFlow.Expression)
	lexer := parser.NewExprLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	p := parser.NewExprParser(stream)
	antlr.ParseTreeWalkerDefault.Walk(e.exprListener, p.Start())
	job := e.exprListener.Pop()
	return job.Exec(context.Background(), workFlow.Para)
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

func (e *Executor) andJob(left, right Job) Job {
	return NewSerialJob(left, right)
}

func (e *Executor) orJob(left, right Job) Job {
	return NewParallelJob(SplitFactory(e.separate), MergeFactory(e.separate), left, right)
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

// Total duration of requests in seconds.
//func initPrometheus(e *Executor, port int) {
//	fieldKeys := []string{"name"}
//
//	e.progress = NewGaugeFrom(prometheus.GaugeOpts{
//		Namespace: PrometheusNamespace,
//		Subsystem: PrometheusSubsystem,
//		Name:      "progress",
//		Help:      "Workflow or Job progress",
//	}, fieldKeys)
//
//	e.latency = NewHistogramFrom(prometheus.HistogramOpts{
//		Namespace: PrometheusNamespace,
//		Subsystem: PrometheusSubsystem,
//		Name:      "latency",
//		Help:      "Workflow or Job exec latency",
//	}, fieldKeys)
//
//	e.success = NewCounterFrom(prometheus.CounterOpts{
//		Namespace: PrometheusNamespace,
//		Subsystem: PrometheusSubsystem,
//		Name:      "success",
//		Help:      "Workflow or Job success count",
//	}, fieldKeys)
//
//	e.failed = NewCounterFrom(prometheus.CounterOpts{
//		Namespace: PrometheusNamespace,
//		Subsystem: PrometheusSubsystem,
//		Name:      "failed",
//		Help:      "Workflow or Job failed count",
//	}, fieldKeys)
//
//	go func() {
//		m := http.NewServeMux()
//		m.Handle("/metrics", promhttp.Handler())
//		_ = http.ListenAndServe(fmt.Sprintf(":%d", port), m)
//	}()
//}

//func (e *Executor) exec(workFlow *WorkFlow) (interface{}, error) {
//	if workFlow == nil {
//		return nil, ErrorInvalidPara
//	}
//	finalJob := NewSerialJob(nil)
//	// [["a"],["a","b","c"], ["x.y.z"]]
//	for _, jobGroup := range workFlow.Expression {
//		parallelJobs := make([]Job, 0)
//		for _, group := range jobGroup {
//			ids := strings.Split(group, ",")
//			if len(ids) > 1 {
//				serial := NewSerialJob(nil)
//				for _, id := range ids {
//					job, err := e.getJob(id)
//					if err != nil {
//						return nil, err
//					}
//					serial.Append(job)
//				}
//				parallelJobs = append(parallelJobs, serial)
//			} else if len(ids) == 1 {
//				job, err := e.getJob(ids[0])
//				if err != nil {
//					return nil, err
//				}
//				parallelJobs = append(parallelJobs, job)
//			} else {
//				return nil, ErrJobNotFound
//			}
//		}
//		// 如果某个节点的job数量大于1
//		// 说明这个节点可以多个job同时运行
//		parallelJob := NewParallelJob(parallelJobs)
//		finalJob.Append(parallelJob)
//	}
//	return finalJob.Exec(context.Background(), workFlow.Para)
//}
