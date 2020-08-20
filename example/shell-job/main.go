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

func (e *ShellJob) Exec(ctx context.Context, req interface{}) (interface{}, error) {
	cmds, err := core.ExactJobRequests(req)
	if err != nil {
		return nil, err
	}

	var resp []byte
	for _, cmd := range cmds {
		fmt.Printf("shell:%s\n", cmd)
		command := exec.Command("/bin/sh", "-c", cmd)
		resp, err = command.Output()
		if err != nil {
			return nil, err
		}
	}
	return string(resp), nil
}

func main() {
	core.StartClientJob(NewShellJob())
}
