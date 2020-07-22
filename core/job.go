package core

import "context"

type Job interface {
	Exec(ctx context.Context, req interface{}) (resp interface{}, err error)
}
