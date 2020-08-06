package core

import (
	"context"
	"fmt"
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

func ExecWorkFlow(flow *WorkFlow) error {
	if flow == nil {
		return ErrorInvalidPara
	}

	// todo
	// 这个地方不太合适
	// mode 暂时对于多个[]job的执行方式
	//for i := range flow.JobIds {
	//	jobs := make([]Job, 0)
	//	for j := range flow.JobIds[i] {
	//
	//	}
	//
	//
	//}

	return nil
}

func ExecJob(para string, mode string, jobs ...Job) (interface{}, error) {
	if len(jobs) == 0 {
		return nil, ErrorInvalidPara
	}
	mode = strings.ToLower(mode)
	switch mode {
	case "serial":
		return ExecJobSerial(para, jobs...)
	case "parallel":
		return ExecJobParallel(para, jobs...)
	default:
		return ExecJobSerial(para, jobs...)
	}
}

func ExecJobSerial(para string, jobs ...Job) (interface{}, error) {
	ctx := context.Background()
	var req interface{}
	req = para
	for i := range jobs {
		resp, err := jobs[i].Exec(ctx, req)
		if err != nil {
			return nil, err
		}
		req = resp
	}
	return req, nil
}

func ExecJobParallel(para string, jobs ...Job) (interface{}, error) {
	var parallelError ParallelError
	ctx := context.Background()
	arr := make([]interface{}, len(jobs))
	wg := sync.WaitGroup{}
	for i := range jobs {
		wg.Add(1)
		go func(pos int) {
			defer wg.Done()
			resp, err := jobs[i].Exec(ctx, para)
			if err != nil {
				parallelError.RawErrors = append(parallelError.RawErrors, fmt.Errorf("[%d] err:%s", pos, err.Error()))
				return
			}
			arr[pos] = resp
		}(i)
	}
	wg.Wait()
	return arr, parallelError
}
