package core

import (
	"github.com/linger1216/jelly-schedule/utils"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const DefaultConfigFilename = "/etc/config/schedule_config.yaml"

func LoadUserConfig(field string, obj interface{}) error {
	configFilename := kingpin.Flag("conf", "config file name").Short('c').Default(DefaultConfigFilename).String()
	kingpin.Version("0.1.0")
	kingpin.Parse()

	buf, err := readFileContent(*configFilename)
	if err != nil {
		return err
	}

	m := make(map[string]interface{})
	err = yaml.Unmarshal(buf, &m)
	if err != nil {
		return err
	}
	return mapstructure.Decode(m[field], obj)
}

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

func readFileContent(filename string) ([]byte, error) {
	obj, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(obj)
	_ = obj.Close()
	return buf, err
}