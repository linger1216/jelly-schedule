package core

import (
	"github.com/linger1216/go-utils/inout"
	"gopkg.in/yaml.v2"
)

type ScheduleConfig struct {
	Etcd     EtcdConfig
	Postgres PostgresConfig
	Http     HttpConfig
	Executor ExecutorConfig
	Job      JobConfig
}

func LoadScheduleConfig(filename string) (*ScheduleConfig, error) {
	buf, err := inout.ReadFileContent(filename)
	if err != nil {
		return nil, err
	}
	return loadScheduleConfig(buf)
}

func loadScheduleConfig(buf []byte) (*ScheduleConfig, error) {
	ret := &ScheduleConfig{}
	err := yaml.Unmarshal(buf, ret)
	if err != nil {
		return nil, err
	}

	if len(ret.Executor.Separate) == 0 {
		ret.Executor.Separate = ";"
	}
	return ret, nil
}
