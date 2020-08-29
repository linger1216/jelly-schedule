package core

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"net/http"
)

// 把Job封装成一个RPC服务
const (
	JsonRPCPath = `/rpc`
)

type JsonRPCService struct {
	job Job
}

func (j *JsonRPCService) Exec(r *http.Request, arg *string, result *string) error {
	resp, err := j.job.Exec(context.Background(), *arg)
	if err != nil {
		return err
	}
	*result = resp
	return nil
}

type JsonRPCServer struct {
	job   Job
	stats JobDescription
}

func newJsonRPCServer(stats JobDescription, job Job) *JsonRPCServer {
	return &JsonRPCServer{stats: stats, job: job}
}

func (d *JsonRPCServer) Start() error {
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")
	server.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	s := &JsonRPCService{d.job}
	err := server.RegisterService(s, "")
	if err != nil {
		return err
	}
	r := mux.NewRouter()
	r.Handle(JsonRPCPath, server)
	return http.ListenAndServe(fmt.Sprintf(":%d", d.stats.Port), r)
}

func (d *JsonRPCServer) Close() {

}
