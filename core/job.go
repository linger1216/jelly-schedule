package core

import (
	"context"
)

type JobConfig struct {
	Ids  map[string]string `json:"ids" yaml:"ids" `
	Host string            `json:"host" yaml:"host" `
}

type Job interface {
	Name() string
	Exec(ctx context.Context, req string) (resp string, err error)
}
