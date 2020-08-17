package core

import (
	"context"
)

type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

type Middleware func(Endpoint) Endpoint

type Job interface {
	Name() string
	Exec(ctx context.Context, req interface{}, stats Endpoint) (resp interface{}, err error)
}
