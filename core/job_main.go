package core

import (
	"github.com/linger1216/go-utils/config"
	"github.com/linger1216/go-utils/sys"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
)

const DefaultConfigFilename = "/etc/config/schedule_config.yaml"

func LoadUserConfig(field string, obj interface{}) error {
	configFilename := kingpin.Flag("conf", "config file name").Short('c').Default(DefaultConfigFilename).String()
	kingpin.Version("0.1.0")
	kingpin.Parse()
	yaml := config.NewYamlReader(*configFilename)
	return yaml.ScanKey(field, obj)
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
		err = os.Setenv(sys.SERVICE_HOST, config.Job.Host)
		if err != nil {
			panic(err)
		}
	}
	id, _ := config.Job.Ids[job.Name()]
	end := make(chan error)
	etcd := NewEtcd(&config.Etcd)
	NewJobServer(etcd, id, job)
	go sys.InterruptHandler(end)
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
