package core

import (
	"fmt"
	"strings"
)

// Job 执行期间的错误, 用于串行
type JobError struct {
	Name    string
	Code    int
	Message string
}

func NewJobError(name string, message string) *JobError {
	return &JobError{Name: name, Message: message, Code: -1}
}

func (j *JobError) Error() string {
	return fmt.Sprintf("id:%s code:%d message:%s", j.Name, j.Code, j.Message)
}

type Errors []error

func (e Errors) Error() string {
	ret := make([]string, len(e))
	for i := range e {
		ret[i] = e[i].Error()
	}
	return strings.Join(ret, ",")
}
