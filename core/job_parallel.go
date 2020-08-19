package core

import (
	"context"
	"fmt"
	"go.uber.org/atomic"
	"strings"
	"sync"
)

type ParallelJob struct {
	jobs     []Job
	progress *atomic.Int32
}

func NewParallelJob(jobs []Job) *ParallelJob {
	return &ParallelJob{jobs: jobs, progress: atomic.NewInt32(0)}
}

func (s *ParallelJob) Name() string {
	names := make([]string, 0, len(s.jobs))
	for _, v := range s.jobs {
		names = append(names, v.Name())
	}
	return strings.Join(names, "-")
}

func (s *ParallelJob) Progress() int {
	return int(s.progress.Load())
}

func (s *ParallelJob) Exec(ctx context.Context, req interface{}) (interface{}, error) {
	reqs, err := exactParallelRequest(req, len(s.jobs))
	l.Debugf("ParallelJob reqs:%v", reqs)
	if err != nil {
		return nil, err
	}

	var rawErrors Errors
	arr := make([]interface{}, len(s.jobs))
	wg := sync.WaitGroup{}
	for i := range s.jobs {
		wg.Add(1)
		go func(pos int) {
			defer wg.Done()
			defer s.progress.Add(int32(100 / len(s.jobs)))
			resp, err := s.jobs[i].Exec(ctx, reqs[i])
			if err != nil {
				rawErrors = append(rawErrors, fmt.Errorf("[%d] err:%s", pos, err.Error()))
				return
			}
			arr[pos] = resp
		}(i)
	}
	wg.Wait()
	s.progress.CAS(int32(s.Progress()), 100)

	if len(rawErrors) > 0 {
		return nil, rawErrors
	}
	return arr, nil
}
