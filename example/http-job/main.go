package main

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/linger1216/jelly-schedule/core"
	"io/ioutil"
	"net/http"
	"strings"
)

import _ "net/http/pprof"

type HttpRequest struct {
	Url    string `json:"url" yaml:"url" `
	Method string `json:"method" yaml:"method" `
	Body   string `json:"body" yaml:"body" `
}

type HttpJob struct {
}

func NewHttpJob() *HttpJob {
	return &HttpJob{}
}

func (e *HttpJob) Name() string {
	return "HttpJob"
}

func (e *HttpJob) Exec(ctx context.Context, req string) (string, error) {
	reqs, err := core.UnMarshalJobRequests(req, ";")
	for i := range reqs {
		for _, arr := range reqs[i].Values {
			for _, cmd := range arr {
				httpRequest := &HttpRequest{}
				err = jsoniter.ConfigFastest.Unmarshal([]byte(cmd), httpRequest)
				if err != nil {
					return "", err
				}
				_, err = doHttpRequest(httpRequest)
				if err != nil {
					return "", err
				}
			}
		}
	}
	return core.GenJobRequestStringByMeta(";", core.NewJobRequestByMeta(reqs...))
}

func doHttpRequest(request *HttpRequest) ([]byte, error) {
	req, err := http.NewRequest(strings.ToUpper(request.Method), request.Url, strings.NewReader(request.Body))
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("err:%s", resp.Status)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func main() {
	core.StartClientJob(NewHttpJob())
}
