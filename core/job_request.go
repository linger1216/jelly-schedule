package core

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"github.com/linger1216/jelly-schedule/utils"
	"strings"
)

type SplitFunc func(sep string, paras string) ([]string, error)
type MergeFunc func(sep string, paras ...string) (string, error)

const (
	EmptyJobRequest = "{}"
)

func _splitStrings(sep string, paras string) ([]string, error) {
	arr := strings.Split(paras, sep)
	ret := make([]string, len(arr))
	for i := range arr {
		ret[i] = arr[i]
	}
	return ret, nil
}

func _mergeStrings(sep string, paras ...string) (string, error) {
	var buf bytes.Buffer
	for i := range paras {
		buf.WriteString(paras[i])
		if i < len(paras)-1 {
			buf.WriteString(sep)
		}
	}
	return buf.String(), nil
}

func _mergeJobRequests(sep string, paras ...string) (string, error) {
	ret := NewJobRequest()
	for i := range paras {
		jobRequest := &JobRequest{}
		if err := jsoniter.ConfigFastest.UnmarshalFromString(paras[i], jobRequest); err != nil {
			return "", err
		}

		// copy values
		for k := range jobRequest.Values {
			if _, ok := ret.Values[k]; ok {
				ret.Values[k] = append(ret.Values[k], jobRequest.Values[k]...)
			} else {
				ret.Values[k] = jobRequest.Values[k]
			}
		}

		// copy meta
		for k, v := range jobRequest.Meta {
			ret.Meta[k] = v
		}
	}
	buf, err := MarshalJobRequests(sep, ret)
	if err != nil {
		return "", err
	}
	return buf, nil
}

/*
Meta 每个request负责解释
Values 呈现给job的值域
Pattern 值域表达式, 负责填充值域

//Values  []string               `json:"values,omitempty"`
*/
type JobRequest struct {
	Meta map[string]interface{} `json:"meta,omitempty"`
	//Values  []string               `json:"values,omitempty"`
	Values  map[string][]string `json:"values,omitempty"`
	Pattern string              `json:"pattern,omitempty"`
	group   int
}

func NewJobRequest() *JobRequest {
	return &JobRequest{Meta: make(map[string]interface{}), Values: make(map[string][]string)}
}

type JobResponse JobRequest

func (j *JobRequest) gen() error {
	if len(j.Pattern) == 0 {
		return nil
	}
	p, err := ParsePattern(j.Pattern)
	if err != nil {
		return err
	}
	j.Values = p.Map(defaultKeyGen)
	j.Pattern = ""
	return nil
}

func (j *JobRequest) split(n int) []*JobRequest {
	total := len(j.Values)
	pages := utils.SplitPage(int64(total), n)
	ret := make([]*JobRequest, 0, n)
	for _, page := range pages {
		req := &JobRequest{
			Meta:   j.Meta,
			Values: make(map[string][]string),
		}
		// random get k,v
		for i := page.Start; i < page.End; i++ {
			var key string
			var val []string
			for k, v := range j.Values {
				key = k
				val = v
				break
			}
			req.Values[key] = val
			delete(j.Values, key)
		}
		ret = append(ret, req)
	}
	if len(ret) == 0 {
		return nil
	}
	return ret
}

func NewJobRequestByKey(key string, src *JobRequest) *JobRequest {
	req := NewJobRequest()
	req.Values[key] = src.Values[key]
	for k, v := range src.Meta {
		req.Meta[k] = v
	}
	return req
}

func NewJobRequestByMeta(src ...*JobRequest) *JobRequest {
	req := NewJobRequest()
	for _, one := range src {
		for k, v := range one.Meta {
			req.Meta[k] = v
		}
	}
	return req
}

func GenJobRequestStringByMeta(sep string, src ...*JobRequest) (string, error) {
	req := NewJobRequest()
	for _, one := range src {
		for k, v := range one.Meta {
			req.Meta[k] = v
		}
	}
	return MarshalJobRequests(sep, req)
}

func MarshalJobRequests(sep string, reqs ...*JobRequest) (string, error) {
	size := len(reqs)
	if size == 0 {
		return "", nil
	}
	paras := make([]string, size)
	for i := range reqs {
		v, err := jsoniter.ConfigFastest.Marshal(reqs[i])
		if err != nil {
			return "", err
		}
		paras[i] = string(v)
	}
	return _mergeStrings(sep, paras...)
}

func UnMarshalJobRequests(req, sep string) ([]*JobRequest, error) {
	paras := strings.Split(req, sep)
	ret := make([]*JobRequest, 0, len(paras))
	for i := range paras {
		jobRequest := NewJobRequest()
		if err := jsoniter.ConfigFastest.UnmarshalFromString(paras[i], jobRequest); err == nil {
			ret = append(ret, jobRequest)
		} else {
			return nil, err
		}
	}
	return ret, nil
}
