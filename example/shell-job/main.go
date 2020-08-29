package main

import (
	"context"
	"fmt"
	"github.com/linger1216/jelly-schedule/core"
	"github.com/linger1216/jelly-schedule/utils"
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
	cmds, err := utils.ExactStringArrayRequests(req, ";")
	if err != nil {
		return "", err
	}

	var resp []byte
	for _, cmd := range cmds {
		fmt.Printf("shell:%s\n", cmd)
		command := exec.Command("/bin/sh", "-c", cmd)
		resp, err = command.Output()
		if err != nil {
			return "", err
		}
	}
	return string(resp), nil
}

func main() {
	core.StartClientJob(NewShellJob())
}
