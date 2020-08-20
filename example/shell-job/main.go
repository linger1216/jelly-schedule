package main

import (
	"context"
	"fmt"
	"github.com/linger1216/jelly-schedule/core"
	"github.com/linger1216/jelly-schedule/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/exec"
)

import _ "net/http/pprof"

const DefaultConfigFilename = "/etc/config/schedule_config.yaml"

var (
	configFilename = kingpin.Flag("conf", "config file name").Short('c').Default(DefaultConfigFilename).String()
)

type ShellJob struct {
}

func NewShellJob() *ShellJob {
	return &ShellJob{}
}

func (e *ShellJob) Name() string {
	return "ShellJob"
}

func (e *ShellJob) Exec(ctx context.Context, req interface{}) (interface{}, error) {
	cmds, err := utils.ExactJobRequests(req)
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

func init() {
	kingpin.Version("0.1.0")
	kingpin.Parse()
}

func main() {
	config, err := core.LoadScheduleConfig(*configFilename)
	if err != nil {
		panic(err)
	}
	if len(config.Job.Host) > 0 {
		err = os.Setenv("SERVICE_HOST", config.Job.Host)
		if err != nil {
			panic(err)
		}
	}
	end := make(chan error)
	etcd := core.NewEtcd(&config.Etcd)
	core.NewJobServer(etcd, NewShellJob())
	go utils.InterruptHandler(end)
	<-end
}
