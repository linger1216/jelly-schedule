package core

import (
	"bytes"
	"context"
	"fmt"
	"github.com/apcera/termtables"
	"github.com/gorilla/rpc/v2/json"
	jsoniter "github.com/json-iterator/go"
	"net/http"
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
	table := termtables.CreateTable()
	table.AddHeaders("Field", "Value")
	table.AddRow("Name", w.Name)
	table.AddRow("Host", w.Host)
	table.AddRow("Port", w.Port)
	table.AddRow("ServicePath", w.ServicePath)
	table.AddRow("JobPath", w.JobPath)
	return table.Render()
}

func (w JobDescription) ToJob() Job {
	return NewDefaultJob(&w)
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

// executor从workflow中得到了job的id
// 利用这个类, 封装成一个Job接口
type DefaultJob struct {
	info *JobDescription
}

func NewDefaultJob(info *JobDescription) *DefaultJob {
	return &DefaultJob{info: info}
}

func (e *DefaultJob) Exec(ctx context.Context, req interface{}) (interface{}, error) {
	message, err := json.EncodeClientRequest("JsonRPCService.exec", req)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("http://%s:%d/%s", e.info.Host, e.info.Port, e.info.ServicePath)
	l.Debugf("%s rpc invoke %s", e.Name(), uri)
	resp, err := http.Post(uri, "application/json", bytes.NewReader(message))
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	reply := ""
	err = json.DecodeClientResponse(resp.Body, &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (e *DefaultJob) Name() string {
	return e.info.Name
}

func (e *DefaultJob) Progress() int {
	// todo
	// 这里还没有设计好
	return 100
}
