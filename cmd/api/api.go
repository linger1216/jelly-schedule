package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/linger1216/jelly-schedule/core"
	"github.com/linger1216/jelly-schedule/utils"
	"github.com/linger1216/jelly-schedule/utils/postgres"
	"gopkg.in/alecthomas/kingpin.v2"

	_ "net/http/pprof"
)

// todo

const (
	DefaultAPIPort = 35744
)

// 172.3.0.122:2379
var (
	etcd            = kingpin.Flag("etcd", "etcd address").Required().String()
	postgresAddress = kingpin.Flag("postgres", "postgres address").Required().String()
	port            = kingpin.Flag("port", "api port").Default("0").Int()
	debugPort       = kingpin.Flag("debugPort", "debug port").Default("0").Int()
)

func init() {
	kingpin.Version("0.1.0")
	kingpin.Parse()
}

func main() {
	if *debugPort > 0 {
		go func() {
			fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", *debugPort), nil))
		}()
	}

	end := make(chan error)

	pg := postgres.NewPostgres(&postgres.PostgresConfig{
		Url: *postgresAddress,
	})

	etcd, err := core.NewEtcd([]string{*etcd}, time.Duration(core.TTL)*time.Second)
	if err != nil {
		panic(err)
	}

	api := core.NewScheduleAPI(etcd, pg)

	if *port == 0 {
		p, err := utils.GetFreePort()
		if err != nil {
			p = DefaultAPIPort
		}
		*port = p
	}

	go func() {
		err := api.Start(*port)
		if err != nil {
			panic(err)
		}
	}()
	go interruptHandler(end)

	fmt.Printf("schedule api run :%d", *port)
	<-end

	_ = etcd.Close()
	_ = pg.Close()
}

func interruptHandler(errc chan<- error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	terminateError := fmt.Errorf("%s", <-c)
	errc <- terminateError
}
