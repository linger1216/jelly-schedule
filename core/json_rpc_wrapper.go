package core

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/linger1216/jelly-schedule/utils"
	"github.com/linger1216/jelly-schedule/utils/snowflake"
	"net/http"
)

const (
	JsonRPCPath = `/rpc`
)

type Request interface{}
type Response interface{}

type DefaultJsonRPCService struct {
	jon Job
}

func (d *DefaultJsonRPCService) Exec(r *http.Request, arg *Request, result *Response) error {
	resp, err := d.jon.Exec(context.Background(), *arg)
	if err != nil {
		return err
	}
	*result = resp
	return nil
}

type StartJsonRPCServerResponse struct {
	Id   string
	Name string
	Host string
	Port int
	Path string
}

type DefaultJsonRPCServer struct {
	id   string
	name string
}

func NewDefaultJsonRPCServer(job Job) *DefaultJsonRPCServer {
	return &DefaultJsonRPCServer{name: job.Name(), id: snowflake.Generate()}
}

func (d *DefaultJsonRPCServer) Start() (*StartJsonRPCServerResponse, error) {
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")
	server.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	s := &DefaultJsonRPCService{}
	err := server.RegisterService(s, d.id)
	if err != nil {
		return nil, err
	}
	r := mux.NewRouter()
	r.Handle(JsonRPCPath, server)

	port, err := utils.GetFreePort()
	if err != nil {
		return nil, err
	}
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		return nil, err
	}
	return &StartJsonRPCServerResponse{
		Id:   d.id,
		Name: d.name,
		Host: utils.GetHost(),
		Port: port,
		Path: JsonRPCPath,
	}, nil
}

func (d *DefaultJsonRPCServer) Close() {

}
