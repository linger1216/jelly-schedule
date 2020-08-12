package core

const (
	StateAvaiable  = "Avaiable"
	StateExecuting = "Executing"
	StateFinish    = "Finish"
)

type WorkFlow struct {
	Id          string     `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	JobIds      [][]string `json:"jobIds,omitempty"`
	Cron        string     `json:"cron,omitempty"`
	Para        string     `json:"para"`
	State       string     `json:"state,omitempty"`
	CreateTime  int64      `json:"createTime,omitempty"`
	UpdateTime  int64      `json:"updateTime,omitempty"`
}
