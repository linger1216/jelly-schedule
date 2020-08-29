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
	mergeFn  MergeFunc
	splitFn  SplitFunc
}

func NewParallelJob(splitFn SplitFunc, mergeFn MergeFunc, jobs ...Job) *ParallelJob {
	return &ParallelJob{splitFn: splitFn, mergeFn: mergeFn, jobs: jobs, progress: atomic.NewInt32(0)}
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

func (s *ParallelJob) Exec(ctx context.Context, req string) (string, error) {

	size := len(s.jobs)
	reqs := s.splitFn(req)
	if len(reqs) != size {
		l.Warnf("ParallelJob actural para size:%d, job:%d", len(reqs), size)
	}

	var rawErrors Errors
	paras := make([]string, len(s.jobs))
	wg := sync.WaitGroup{}
	for i := range s.jobs {
		wg.Add(1)
		go func(pos int) {
			defer wg.Done()
			defer s.progress.Add(int32(100 / len(s.jobs)))
			var para string
			if pos < len(reqs) {
				para = reqs[pos]
			}
			resp, err := s.jobs[pos].Exec(ctx, para)
			if err != nil {
				rawErrors = append(rawErrors, fmt.Errorf("[%d] err:%s", pos, err.Error()))
				return
			}
			paras[pos] = resp
		}(i)
	}
	wg.Wait()
	s.progress.CAS(int32(s.Progress()), 100)

	if len(rawErrors) > 0 {
		return "", rawErrors
	}

	// merge parameters
	return s.mergeFn(paras...), nil
}
