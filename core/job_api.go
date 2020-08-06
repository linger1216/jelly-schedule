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

type getJobListRequest struct{}

type getJobListResponse struct {
	JobStats []*JobStats `json:"jobStats"`
}

func decodeGetJobListRequest(r *http.Request) (interface{}, error) {
	pathParams := mux.Vars(r)
	_ = pathParams
	queryParams := r.URL.Query()
	_ = queryParams
	return &getJobListRequest{}, nil
}

func (w *jobAPI) getJobList(ctx context.Context, req interface{}) (interface{}, error) {
	request, ok := req.(*getJobListRequest)
	if !ok {
		return nil, ErrorBadRequest
	}

	_ = request

	_, v, err := w.etcd.GetWithPrefixKey(ctx, JobPrefix)
	if err != nil {
		return nil, err
	}

	resp := &getJobListResponse{}
	for i := range v {
		stats := &JobStats{}
		err = jsoniter.ConfigFastest.Unmarshal([]byte(v[i]), stats)
		if err != nil {
			return nil, err
		}
		resp.JobStats = append(resp.JobStats, stats)
	}

	return resp, nil
}

type getJobRequest struct {
	ids []string
}

type getJobResponse struct {
	JobStats []*JobStats `json:"jobStats"`
}

func decodeGetJobRequest(r *http.Request) (interface{}, error) {
	pathParams := mux.Vars(r)
	_ = pathParams
	queryParams := r.URL.Query()
	_ = queryParams
	return &getJobRequest{
		ids: strings.Split(pathParams["ids"], ","),
	}, nil
}

func (w *jobAPI) getJob(ctx context.Context, req interface{}) (interface{}, error) {
	request, ok := req.(*getJobRequest)
	if !ok {
		return nil, ErrorBadRequest
	}
	resp := &getJobResponse{}
	for i := range request.ids {
		v, err := w.etcd.Get(ctx, JobPrefix+"/"+request.ids[i])
		if err != nil {
			return nil, err
		}
		stats := &JobStats{}
		err = jsoniter.ConfigFastest.Unmarshal(v, stats)
		if err != nil {
			return nil, err
		}
		resp.JobStats = append(resp.JobStats, stats)
	}
	return resp, nil
}
