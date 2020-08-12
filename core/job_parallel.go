package core

import (
	"context"
	"fmt"
	"go.uber.org/atomic"
	"strings"
	"sync"
)

type ParallelError struct {
	RawErrors []error
	Final     error
}

func (e ParallelError) Error() string {
	var suffix string
	if len(e.RawErrors) > 1 {
		a := make([]string, len(e.RawErrors)-1)
		for i := 0; i < len(e.RawErrors)-1; i++ { // last one is Final
			a[i] = e.RawErrors[i].Error()
		}
		suffix = fmt.Sprintf(" (previously: %s)", strings.Join(a, "; "))
	}
	return fmt.Sprintf("%v%s", e.Final, suffix)
}

func (e ParallelError) Empty() bool {
	return len(e.RawErrors) == 0
}

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
	var parallelError ParallelError
	arr := make([]interface{}, len(s.jobs))
	wg := sync.WaitGroup{}
	for i := range s.jobs {
		wg.Add(1)
		go func(pos int) {
			defer wg.Done()
			defer s.progress.Add(int32(100 / len(s.jobs)))
			resp, err := s.jobs[i].Exec(ctx, req)
			if err != nil {
				parallelError.RawErrors = append(parallelError.RawErrors, fmt.Errorf("[%d] err:%s", pos, err.Error()))
				return
			}
			arr[pos] = resp
		}(i)
	}
	wg.Wait()
	s.progress.CAS(int32(s.Progress()), 100)

	if !parallelError.Empty() {
		return nil, parallelError
	}
	return arr, nil
}
