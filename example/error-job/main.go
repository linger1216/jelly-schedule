package main

import (
	"context"
	"fmt"
	"github.com/linger1216/jelly-schedule/core"
	"time"
)

import _ "net/http/pprof"

type ErrorJob struct {
}

func NewErrorJob() *ErrorJob {
	return &ErrorJob{}
}

func (e *ErrorJob) Name() string {
	return "ErrorJob"
}

func (e *ErrorJob) Exec(ctx context.Context, req interface{}) (resp interface{}, err error) {
	fmt.Printf("ErrorJob:%d\n", time.Now().Unix())
	return nil, fmt.Errorf("fake error")
}

func main() {
	core.StartClientJob(NewErrorJob())
}
