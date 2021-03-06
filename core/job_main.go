package core

import (
	"github.com/linger1216/go-utils/config"
	"github.com/linger1216/go-utils/sys"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
)

const DefaultConfigFilename = "/etc/config/schedule_config.yaml"

var (
	configFilename *string
	yamlReader     *config.YamlReader
)

func init() {
	configFilename = kingpin.Flag("conf", "config file name").Short('c').Default(DefaultConfigFilename).String()
	kingpin.Version("0.1.0")
	kingpin.Parse()
	yamlReader = config.NewYamlReader(*configFilename)
}

func LoadUserConfig(field string, obj interface{}) error {
	return yamlReader.ScanKey(field, obj)
}

func StartClientJob(job Job) {
	configFilename := kingpin.Flag("conf", "config file name").Short('c').Default(DefaultConfigFilename).String()

	kingpin.Version("0.1.0")
	kingpin.Parse()

	conf, err := LoadScheduleConfig(*configFilename)
	if err != nil {
		panic(err)
	}

	if len(conf.Job.Host) > 0 {
		err = os.Setenv(sys.SERVICE_HOST, conf.Job.Host)
		if err != nil {
			panic(err)
		}
	}
	id, _ := conf.Job.Ids[job.Name()]
	end := make(chan error)
	etcd := NewEtcd(&conf.Etcd)
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
