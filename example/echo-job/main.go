package main

import (
	"context"
	"fmt"
	"github.com/linger1216/jelly-schedule/core"
	"time"
)

import _ "net/http/pprof"

type EchoJob struct {
}

func NewEchoJob() *EchoJob {
	return &EchoJob{}
}

func (e *EchoJob) Name() string {
	return "EchoJob"
}

func (e *EchoJob) Exec(ctx context.Context, req string) (resp string, err error) {
	fmt.Printf("echo:%s (%d)\n", req, time.Now().Unix())
	return req, nil
}

func main() {
	core.StartClientJob(NewEchoJob())
}
