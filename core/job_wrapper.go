package core

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gorilla/rpc/json"
	"net/http"
)

const (
	RemoteServerMethod = `JsonRPCService.Exec`
)

// executor从workflow中得到了job的id
// 利用这个类, 封装成一个Job接口
type WrapperJob struct {
	info *JobDescription
}

func NewWrapperJob(info *JobDescription) *WrapperJob {
	return &WrapperJob{info: info}
}

func (e *WrapperJob) Exec(ctx context.Context, req string) (string, error) {
	message, err := json.EncodeClientRequest(RemoteServerMethod, req)
	if err != nil {
		return "", err
	}

	uri := fmt.Sprintf("http://%s:%d/%s", e.info.Host, e.info.Port, e.info.ServicePath)
	//l.Debugf("%s rpc invoke %s", e.Name(), uri)
	resp, err := http.Post(uri, "application/json", bytes.NewReader(message))
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	reply := ""
	err = json.DecodeClientResponse(resp.Body, &reply)
	if err != nil {
		return "", err
	}
	return reply, nil
}

func (e *WrapperJob) Name() string {
	return e.info.Name
}
