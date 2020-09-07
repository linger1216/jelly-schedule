package core

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/atomic"
	"strings"
)

type SerialJob struct {
	sep      string
	jobs     []Job
	progress *atomic.Int32
}

func NewSerialJob(sep string, jobs ...Job) *SerialJob {
	return &SerialJob{sep: sep, jobs: jobs, progress: atomic.NewInt32(0)}
}

func (s *SerialJob) Append(job Job) {
	s.jobs = append(s.jobs, job)
}

func (s *SerialJob) Name() string {
	names := make([]string, 0, len(s.jobs))
	for _, v := range s.jobs {
		names = append(names, v.Name())
	}
	return strings.Join(names, "->")
}

func (s *SerialJob) Progress() int {
	return int(s.progress.Load())
}

func (s *SerialJob) Exec(ctx context.Context, req string) (string, error) {
	// 串行任务收到的永远是一个request
	// 并行任务收到的可能是一个request或者[]request
	arg := req
	_MOD(_SerialJob).With(_Job, s.Name()).Debugf("exec req:%s", req)
	for i := range s.jobs {
		// 后面有n个并发任务, 检查参数是不是要分割
		n := 1
		switch x := s.jobs[i].(type) {
		case *SerialJob:
		case *AlternateJob:
		case *ParallelJob:
			n = len(x.jobs)
		}

		jobRequest := NewJobRequest()
		if err := jsoniter.ConfigFastest.UnmarshalFromString(arg, jobRequest); err != nil {
			return "", err
		}

		if err := jobRequest.gen(); err != nil {
			return "", err
		}

		jobRequestsStr, err := marshalJobRequests(s.sep, jobRequest.split(n)...)
		if err != nil {
			return "", err
		}
		resp, err := s.jobs[i].Exec(ctx, jobRequestsStr)
		if err != nil {
			return "", err
		}
		s.progress.Add(int32(100 / len(s.jobs)))
		arg = resp
	}
	s.progress.CAS(int32(s.Progress()), 100)
	return arg, nil
}
