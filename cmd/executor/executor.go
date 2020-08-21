package main

import (
	"github.com/linger1216/jelly-schedule/core"
	"gopkg.in/alecthomas/kingpin.v2"
	_ "net/http/pprof"
)

const DefaultConfigFilename = "/etc/config/schedule_config.yaml"

var (
	configFilename = kingpin.Flag("conf", "config file name").Short('c').Default(DefaultConfigFilename).String()
)

func init() {
	kingpin.Version("0.1.0")
	kingpin.Parse()
}

func main() {
	config, err := core.LoadScheduleConfig(*configFilename)
	if err != nil {
		panic(err)
	}
	end := make(chan error)
	pg := core.NewPostgres(&config.Postgres)
	etcd := core.NewEtcd(&config.Etcd)
	exec := core.NewExecutor(etcd, pg, config.Executor)
	_ = exec
	<-end
	_ = etcd.Close()
	_ = pg.Close()
}
