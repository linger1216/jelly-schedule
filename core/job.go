package core

import (
	"context"
)

type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
type Middleware func(Endpoint) Endpoint

type JobConfig struct {
	Host string `json:"host" yaml:"host" `
}

type Job interface {
	Name() string
	Exec(ctx context.Context, req interface{}) (resp interface{}, err error)
}
