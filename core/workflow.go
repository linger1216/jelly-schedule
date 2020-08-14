package core

import "go.uber.org/atomic"

type WorkFlow struct {
	Id          string     `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	JobIds      [][]string `json:"jobIds,omitempty"`
	Cron        string     `json:"cron,omitempty"`
	Para        string     `json:"para"`
	// 执行几次结束
	ExecuteLimit int64 `json:"executeLimit" yaml:"executeLimit" `
	// 碰到错误的方式
	ErrorPolicy string `json:"errorPolicy" yaml:"errorPolicy"`
	// 可以指定由哪个执行器执行
	BelongExecutor string `json:"belongExecutor" yaml:"belongExecutor" `
	State          string `json:"state,omitempty"`
	CreateTime     int64  `json:"createTime,omitempty"`
	UpdateTime     int64  `json:"updateTime,omitempty"`
}

type WorkFlowStats struct {
	Id                   string
	SuccessExecuteCount  atomic.Int32
	RetryExecuteCount    atomic.Int32
	FailedExecuteCount   atomic.Int32
	LastExecuteDuration  atomic.Int32
	TotalExecuteDuration atomic.Int32
}
