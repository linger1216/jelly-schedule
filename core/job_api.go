package core

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"
)

type jobAPI struct {
	etcd *Etcd
}

func NewJobAPI(etcd *Etcd) *jobAPI {
	return &jobAPI{etcd: etcd}
}

type ListJobRequest struct{}

type ListJobResponse struct {
	JobStats []*JobInfo `json:"jobStats"`
}

func NewListJobResponse() *ListJobResponse {
	return &ListJobResponse{JobStats: make([]*JobInfo, 0)}
}

func decodeGetJobListRequest(r *http.Request) (interface{}, error) {
	pathParams := mux.Vars(r)
	_ = pathParams
	queryParams := r.URL.Query()
	_ = queryParams
	return &ListJobRequest{}, nil
}

func (w *jobAPI) listJob(ctx context.Context, req interface{}) (interface{}, error) {
	_, ok := req.(*ListJobRequest)
	if !ok {
		return nil, ErrBadRequest
	}

	_, v, err := w.etcd.GetWithPrefixKey(ctx, JobPrefix)
	if err != nil {
		return nil, err
	}

	resp := NewListJobResponse()
	if len(v) == 0 {
		return resp, nil
	}

	for i := range v {
		stats := &JobInfo{}
		err = jsoniter.ConfigFastest.Unmarshal([]byte(v[i]), stats)
		if err != nil {
			return nil, err
		}
		resp.JobStats = append(resp.JobStats, stats)
	}
	return resp, nil
}

type GetJobRequest struct {
	ids []string
}

type GetJobResponse struct {
	JobStats []*JobInfo `json:"jobStats"`
}

func NewGetJobResponse() *GetJobResponse {
	return &GetJobResponse{JobStats: make([]*JobInfo, 0)}
}

func decodeGetJobRequest(r *http.Request) (interface{}, error) {
	pathParams := mux.Vars(r)
	_ = pathParams
	queryParams := r.URL.Query()
	_ = queryParams
	return &GetJobRequest{
		ids: strings.Split(pathParams["ids"], ","),
	}, nil
}

func (w *jobAPI) getJob(ctx context.Context, req interface{}) (interface{}, error) {
	request, ok := req.(*GetJobRequest)
	if !ok {
		return nil, ErrBadRequest
	}
	resp := NewGetJobResponse()
	for i := range request.ids {
		v, err := w.etcd.Get(ctx, JobKey(request.ids[i]))
		if err != nil {
			return nil, err
		}
		if len(v) == 0 {
			continue
		}

		stats := &JobInfo{}
		err = jsoniter.ConfigFastest.Unmarshal(v, stats)
		if err != nil {
			return nil, err
		}
		resp.JobStats = append(resp.JobStats, stats)
	}
	return resp, nil
}

func encodeHTTPJobResponse(w http.ResponseWriter, response interface{}) error {
	encoder := jsoniter.ConfigFastest.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	switch x := response.(type) {
	case *ListJobResponse:
		return encoder.Encode(x.JobStats)
	case *GetJobResponse:
		return encoder.Encode(x.JobStats)
	}
	return encoder.Encode(response)
}
