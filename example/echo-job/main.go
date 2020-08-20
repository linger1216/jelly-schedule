package main

import (
	"context"
	"fmt"
	"github.com/linger1216/jelly-schedule/core"
	"github.com/linger1216/jelly-schedule/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

import _ "net/http/pprof"

const DefaultConfigFilename = "/etc/config/schedule_config.yaml"

var (
	configFilename = kingpin.Flag("conf", "config file name").Short('c').Default(DefaultConfigFilename).String()
)

type EchoJob struct {
}

func NewEchoJob() *EchoJob {
	return &EchoJob{}
}

func (e *EchoJob) Name() string {
	return "EchoJob"
}

func (e *EchoJob) Exec(ctx context.Context, req interface{}) (resp interface{}, err error) {
	cmds, err := utils.ExactJobRequests(req)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("echo:%s\n", strings.Join(cmds, ",")), nil
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
	core.NewJobServer(etcd, NewEchoJob())
	go interruptHandler(end)
	<-end
}

func interruptHandler(errc chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	terminateError := fmt.Errorf("%s", <-c)
	errc <- terminateError
}
