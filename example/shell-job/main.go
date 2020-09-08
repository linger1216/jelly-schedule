package main

import (
	"context"
	"fmt"
	"github.com/linger1216/jelly-schedule/core"
	"os/exec"
)

import _ "net/http/pprof"

type ShellJob struct {
}

func NewShellJob() *ShellJob {
	return &ShellJob{}
}

func (e *ShellJob) Name() string {
	return "ShellJob"
}

func (e *ShellJob) Exec(ctx context.Context, req string) (string, error) {
	reqs, err := core.UnMarshalJobRequests(req, ";")
	for i := range reqs {
		for _, arr := range reqs[i].Values {
			for _, cmd := range arr {
				fmt.Printf("shell:%s\n", cmd)
				command := exec.Command("/bin/sh", "-c", cmd)
				_, err = command.Output()
				if err != nil {
					return "", err
				}
			}
		}
	}
	return core.GenJobRequestStringByMeta(";", core.NewJobRequestByMeta(reqs...))
}

func main() {
	core.StartClientJob(NewShellJob())
}
