package main

import (
	"context"
	"fmt"
	"github.com/linger1216/jelly-schedule/core"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

import _ "net/http/pprof"

var (
	etcd  = kingpin.Flag("etcd", "etcd address").Required().String()
	debug = kingpin.Flag("debug", "debug debug").Default("0").Int()
)

type EchoJob struct {
}

func NewEchoJob() *EchoJob {
	return &EchoJob{}
}

func (e *EchoJob) Name() string {
	return "EchoJob"
}

func (e *EchoJob) Progress() int {
	return 100
}

func (e *EchoJob) Exec(ctx context.Context, req interface{}) (resp interface{}, err error) {
	fmt.Printf("echo:%s\n", req.(string))
	return "ok", nil
}


func init() {
	kingpin.Version("0.1.0")
	kingpin.Parse()
}

func main() {
	if *debug > 0 {
		go func() {
			log.Println(http.ListenAndServe(fmt.Sprintf(":%d", *debug), nil))
		}()
	}

	end := make(chan error)
	etcd, err := core.NewEtcd([]string{*etcd}, time.Duration(core.TTL)*time.Second)
	if err != nil {
		panic(err)
	}
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
