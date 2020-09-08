package core

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/atomic"
	"strings"
)

type LoopJob struct {
	sep      string
	mergeFn  MergeFunc
	jobs     []Job
	progress *atomic.Int32
}

func NewLoopJob(sep string, mergeFn MergeFunc, jobs ...Job) *LoopJob {
	return &LoopJob{sep: sep, mergeFn: mergeFn, jobs: jobs, progress: atomic.NewInt32(0)}
}

func (s *LoopJob) Append(job Job) {
	s.jobs = append(s.jobs, job)
}

func (s *LoopJob) Name() string {
	names := make([]string, 0, len(s.jobs))
	for _, v := range s.jobs {
		names = append(names, v.Name())
	}
	return strings.Join(names, "=>")
}

func (s *LoopJob) Progress() int {
	return int(s.progress.Load())
}

func (s *LoopJob) Exec(ctx context.Context, req string) (string, error) {

	rawRequest := NewJobRequest()
	if err := jsoniter.ConfigFastest.UnmarshalFromString(req, rawRequest); err != nil {
		return "", err
	}
	if err := rawRequest.gen(); err != nil {
		return "", err
	}

	resps := make([]string, 0, len(rawRequest.Values))
	for key := range rawRequest.Values {
		_MOD(_AlternateJob).With(_Job, s.Name()).Debugf("req :%s", req)

		// 产生一个新的request
		oneRequest := NewJobRequestByKey(key, rawRequest)
		oneRequestBuf, err := MarshalJobRequests(s.sep, oneRequest)
		if err != nil {
			return "", err
		}

		arg := oneRequestBuf
		for j := range s.jobs {
			// 这时候任务可能是串/并/交替
			// 但不管是什么, 只传给一个request,
			// 讲道理, 后面跟并行你需要再仔细考虑下, 除非你知道你再干什么?
			// (讲道理作为作者的我, 都有点晕)
			resp, err := s.jobs[j].Exec(ctx, arg)
			if err != nil {
				return "", err
			}
			s.progress.Add(int32(100 / len(s.jobs)))
			arg = resp
		}

		// 保存所有的参数, 理论上每一个都是一个jobRequest
		resps = append(resps, arg)
	}

	s.progress.CAS(int32(s.Progress()), 100)

	// merge parameters
	return s.mergeFn(s.sep, resps...)
}
