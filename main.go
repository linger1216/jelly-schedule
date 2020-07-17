package main

import (
	"fmt"
	"github.com/linger1216/jelly-schedule/core"
	"github.com/linger1216/jelly-schedule/etcdv3"
	"github.com/linger1216/jelly-schedule/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	name = kingpin.Flag("name", "worker name").String()
	etcd = kingpin.Flag("etcd", "etcd address").Default("127.0.0.1:2379").String()
)

func init() {
	kingpin.Version("0.1.0")
	kingpin.Parse()
}

func main() {
	if len(*name) == 0 {
		*name = utils.GetHost()
	}

	end := make(chan error)
	etcd, err := etcdv3.NewEtcd([]string{*etcd}, time.Second)
	if err != nil {
		panic(err)
	}
	core.NewWorker(*name, etcd)
	go interruptHandler(end)
	<-end
}

func interruptHandler(errc chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	terminateError := fmt.Errorf("%s", <-c)
	errc <- terminateError
}
