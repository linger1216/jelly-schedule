package core

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"strings"
)

type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

type Middleware func(Endpoint) Endpoint

type Job interface {
	Name() string
	Exec(ctx context.Context, req interface{}) (resp interface{}, err error)
}

type JobConfig struct {
	Host string `json:"name" yaml:"name" `
}

func exactSerialRequest(req interface{}) interface{} {
	var arg interface{}
	switch x := req.(type) {
	case string:
		arg = x
	case []interface{}:
		if len(x) == 1 {
			arg = x[0]
		} else {
			arg = x
		}
	}
	return arg
}

func exactParallelRequest(req interface{}, size int) ([]interface{}, error) {
	args := make([]interface{}, 0, size)
	switch x := req.(type) {
	case string:
		arr := make([]string, 0)
		_ = jsoniter.ConfigFastest.Unmarshal([]byte(x), &arr)
		if len(arr) == 0 {
			arr = strings.Split(x, ",")
		}
		if len(arr) != size {
			return nil, ErrorJobParaInvalid
		}
		for i := range arr {
			args = append(args, arr[i])
		}
	case []interface{}:
		// 后面是一个单独的任务, 要保证后续的参数正常
		if size > 1 && len(x) != size {
			return nil, ErrorJobParaInvalid
		} else {
			args = x
		}
	}
	return args, nil
}
