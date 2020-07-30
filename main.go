package main

import (
	"fmt"
	"github.com/linger1216/jelly-schedule/core"
	"github.com/linger1216/jelly-schedule/utils"
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
	name = kingpin.Flag("name", "worker name").String()
	etcd = kingpin.Flag("etcd", "etcd address").Default("172.3.0.122:2379").String()
)

func init() {
	kingpin.Version("0.1.0")
	kingpin.Parse()
}

func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	if len(*name) == 0 {
		*name = utils.GetHost()
	}

	end := make(chan error)
	etcd, err := core.NewEtcd([]string{*etcd}, time.Duration(core.TTL)*time.Second)
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
