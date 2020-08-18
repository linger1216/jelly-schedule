package core

type WorkFlow struct {
	Id          string     `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	JobIds      [][]string `json:"jobIds,omitempty"`
	Cron        string     `json:"cron,omitempty"`
	Para        string     `json:"para"`
	// 执行几次结束
	SuccessLimit int64 `json:"successLimit" yaml:"successLimit" `
	// 碰到错误的方式
	FailedLimit int64 `json:"failedLimit" yaml:"failedLimit"`
	// 可以指定由哪个执行器执行
	BelongExecutor string `json:"belongExecutor" yaml:"belongExecutor" `
	State          string `json:"state,omitempty"`
	CreateTime     int64  `json:"createTime,omitempty"`
	UpdateTime     int64  `json:"updateTime,omitempty"`
}
