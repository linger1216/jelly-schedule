package core

import (
	"bytes"
	jsoniter "github.com/json-iterator/go"
	"github.com/linger1216/jelly-schedule/utils"
	"strings"
)

type SplitFunc func(sep string, paras string) ([]string, error)
type MergeFunc func(sep string, paras ...string) (string, error)

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
		ret.Values = append(ret.Values, jobRequest.Values...)
		for k, v := range jobRequest.Meta {
			ret.Meta[k] = v
		}
	}
	buf, err := marshalJobRequests(sep, ret)
	if err != nil {
		return "", err
	}
	return buf, nil
}

/*
Meta 每个request负责解释
Values 呈现给job的值域
Pattern 值域表达式, 负责填充值域
*/
type JobRequest struct {
	Meta    map[string]interface{} `json:"meta,omitempty"`
	Values  []string               `json:"values,omitempty"`
	Pattern string                 `json:"pattern,omitempty"`
	group   int
}

func NewJobRequest() *JobRequest {
	return &JobRequest{Meta: make(map[string]interface{})}
}

type JobResponse JobRequest

func (j *JobRequest) gen() error {

	if len(j.Pattern) == 0 {
		return nil
	}

	splitHolders, arrangeHolders, err := parsePattern(j.Pattern)
	if err != nil {
		return err
	}
	splitPatterns, err := arrangePattern([]string{j.Pattern}, splitHolders)
	if err != nil {
		return err
	}
	for _, v := range splitPatterns {
		if arrangePatterns, err := arrangePattern([]string{v}, arrangeHolders); err == nil {
			j.Values = append(j.Values, arrangePatterns...)
		}
	}

	j.Pattern = ""
	return nil
}

func (j *JobRequest) split(n int) []*JobRequest {
	// todo
	_ = j.group
	total := len(j.Values)
	pages := utils.SplitPage(int64(total), n)
	ret := make([]*JobRequest, 0, n)
	for _, page := range pages {
		req := &JobRequest{
			Meta:   j.Meta,
			Values: j.Values[page.Start:page.End],
		}
		ret = append(ret, req)
	}
	if len(ret) == 0 {
		return nil
	}
	return ret
}

func marshalJobRequests(sep string, reqs ...*JobRequest) (string, error) {
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
