package core

import "context"

type Job interface {
	Name() string
	Exec(ctx context.Context, req interface{}) (resp interface{}, err error)
}
