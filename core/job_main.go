package core

import (
	"github.com/linger1216/jelly-schedule/utils"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

const DefaultConfigFilename = "/etc/config/schedule_config.yaml"

func StartClientJob(job Job) {
	configFilename := kingpin.Flag("conf", "config file name").Short('c').Default(DefaultConfigFilename).String()

	kingpin.Version("0.1.0")
	kingpin.Parse()

	config, err := LoadScheduleConfig(*configFilename)
	if err != nil {
		panic(err)
	}

	if len(config.Job.Host) > 0 {
		err = os.Setenv(utils.SERVICE_HOST, config.Job.Host)
		if err != nil {
			panic(err)
		}
	}
	id, _ := config.Job.Ids[job.Name()]
	end := make(chan error)
	etcd := NewEtcd(&config.Etcd)
	NewJobServer(etcd, id, job)
	go utils.InterruptHandler(end)
	<-end
}
