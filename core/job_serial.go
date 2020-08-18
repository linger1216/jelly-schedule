package core

import (
	"context"
	"go.uber.org/atomic"
	"strings"
)

type SerialJob struct {
	jobs     []Job
	progress *atomic.Int32
}

func NewSerialJob(jobs []Job) *SerialJob {
	return &SerialJob{jobs: jobs, progress: atomic.NewInt32(0)}
}

func (s *SerialJob) Append(job Job) {
	s.jobs = append(s.jobs, job)
}

func (s *SerialJob) Name() string {
	names := make([]string, 0, len(s.jobs))
	for _, v := range s.jobs {
		names = append(names, v.Name())
	}
	return strings.Join(names, "-")
}

func (s *SerialJob) Progress() int {
	return int(s.progress.Load())
}

func (s *SerialJob) Exec(ctx context.Context, req interface{}) (interface{}, error) {
	arg := req
	for i := range s.jobs {
		resp, err := s.jobs[i].Exec(ctx, arg)
		if err != nil {
			//return nil, NewJobError(s.jobs[i].Name(), err.Error())
			return nil, err
		}
		s.progress.Add(int32(100 / len(s.jobs)))
		arg = resp
	}
	s.progress.CAS(int32(s.Progress()), 100)
	return arg, nil
}
