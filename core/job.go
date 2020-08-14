package core

import (
	"context"
	"fmt"
)

type Job interface {
	Name() string
	Progress() int
	Exec(ctx context.Context, req interface{}) (resp interface{}, err error)
}
