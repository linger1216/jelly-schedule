package core

import (
	"github.com/linger1216/jelly-schedule/utils"
	"gopkg.in/yaml.v2"
)

type ScheduleConfig struct {
	Etcd     EtcdConfig
	Postgres PostgresConfig
	Http     HttpConfig
}

func LoadScheduleConfig(filename string) (*ScheduleConfig, error) {
	buf, err := utils.ReadFileContent(filename)
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
	return ret, nil
}
