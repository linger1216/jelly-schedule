package core

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
)

type JobDescription struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	ServicePath string `json:"servicePath"`
	JobPath     string `json:"jobPath"`
}

func (w JobDescription) String() string {
	return fmt.Sprintf("name:%s host:%s port:%d servicePath:%s jobPath:%s",
		w.Name, w.Host, w.Port, w.ServicePath, w.JobPath)
}

func (w JobDescription) ToJob() Job {
	return NewWrapperJob(&w)
}

func MarshalJobDescription(j *JobDescription) ([]byte, error) {
	return jsoniter.ConfigFastest.Marshal(j)
}

func UnMarshalJobDescription(buf []byte) (*JobDescription, error) {
	s := &JobDescription{}
	err := jsoniter.ConfigFastest.Unmarshal(buf, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
