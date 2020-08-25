package main

import (
	"context"
	"fmt"
	"github.com/linger1216/jelly-schedule/core"
	"strings"
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

func (e *EchoJob) Exec(ctx context.Context, req interface{}) (resp interface{}, err error) {
	fmt.Printf("EchoJob:%d\n", time.Now().Unix())
	cmds, err := core.ExactJobRequests(req)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("echo:%s\n", strings.Join(cmds, ",")), nil
}

func main() {
	core.StartClientJob(NewEchoJob())
}
