package main

import (
	"fmt"
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
	etcd := core.NewEtcd(&config.Etcd)

	for i := 0; i < 100; i++ {
		err = etcd.TryLockWithTTL("/locker/0800", 60)
		if err != nil {
			fmt.Printf("try lock:%s\n", err.Error())
		} else {
			fmt.Printf("locker successful\n")
		}
	}
	_ = etcd.Close()
}
