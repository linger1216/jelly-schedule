package core

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
	"github.com/linger1216/go-pipeline/pipe"
)

type workFlowAPI struct {
	db                   *sqlx.DB
	createWorkflowFilter pipe.Filter
}

func NewWorkflowAPI(db *sqlx.DB) *workFlowAPI {
	/*
		api.createWorkflowFilter = pipe.NewStraightPipeline(false, "create asset").
		Append("validCreateAssets", ret.validCreateAssets).
		Append("validCreateAssets", ret.cleanUpsertAssets).
		Append("execUpsertAssets", ret.execUpsertAssets).
		Append("execUpsertAssets", ret.createAssetsResponse)
	*/
	return &workFlowAPI{db: db}
}

type CreateWorkflowRequest struct {
	Workflows []*WorkFlow
}

type CreateWorkflowResponse struct {
	Ids []string
}

func decodeCreateWorkflowRequest(r *http.Request) (interface{}, error) {
	req := &CreateWorkflowRequest{}
	// to support gzip input
	var reader io.ReadCloser
	switch r.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err := gzip.NewReader(r.Body)
		defer reader.Close()
		if err != nil {
			return nil, newApiError(http.StatusBadRequest, "failed to read the gzip content")
		}
	default:
		reader = r.Body
	}

	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, newApiError(http.StatusBadRequest, "cannot read body of http request")
	}
	if len(buf) > 0 {
		if err = jsoniter.ConfigFastest.Unmarshal(buf, &req.Workflows); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, newApiError(http.StatusBadRequest, fmt.Sprintf("request body '%s': cannot parse non-json request body", buf))
		}
	}
	pathParams := mux.Vars(r)
	_ = pathParams
	queryParams := r.URL.Query()
	_ = queryParams
	return req, nil
}

func (w *workFlowAPI) CreateWorkflow(ctx context.Context, req interface{}) (interface{}, error) {
	request, ok := req.(*CreateWorkflowRequest)
	if !ok {
		return nil, ErrorBadRequest
	}

	in := request.Workflows
	if len(in) == 0 {
		return nil, ErrorInvalidPara
	}

	query, args, err := upsertWorkflowSql(in)
	if err != nil {
		return nil, err
	}

	l.Debug(query)
	_, err = w.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	resp := &CreateWorkflowResponse{}
	for i := range in {
		resp.Ids = append(resp.Ids, in[i].Id)
	}

	return resp, nil
}

type GetWorkflowRequest struct {
	Header       int
	Names        []string
	Descriptions []string
	StartTime    int64
	EndTime      int64
	CurrentPage  uint64
	PageSize     uint64
}

type GetWorkflowResponse struct {
	WorkflowStats []string `json:"WorkflowStats"`
}

func decodeGetWorkflowRequest(r *http.Request) (interface{}, error) {
	pathParams := mux.Vars(r)
	_ = pathParams
	queryParams := r.URL.Query()
	_ = queryParams
	return &GetWorkflowRequest{}, nil
}

func (w *workFlowAPI) GetWorkflow(ctx context.Context, req interface{}) (interface{}, error) {
	request, ok := req.(*GetWorkflowRequest)
	if !ok {
		return nil, ErrorBadRequest
	}

	_ = request
	resp := &GetWorkflowResponse{}
	return resp, nil
}

type UpdateWorkflowRequest struct {
	Header       int
	Names        []string
	Descriptions []string
	StartTime    int64
	EndTime      int64
	CurrentPage  uint64
	PageSize     uint64
}

type UpdateWorkflowResponse struct {
	WorkflowStats []string `json:"WorkflowStats"`
}

func decodeUpdateWorkflowRequest(r *http.Request) (interface{}, error) {
	pathParams := mux.Vars(r)
	_ = pathParams
	queryParams := r.URL.Query()
	_ = queryParams
	return &UpdateWorkflowRequest{}, nil
}

func (w *workFlowAPI) UpdateWorkflow(ctx context.Context, req interface{}) (interface{}, error) {
	request, ok := req.(*UpdateWorkflowRequest)
	if !ok {
		return nil, ErrorBadRequest
	}

	_ = request

	resp := &UpdateWorkflowResponse{}

	return resp, nil
}

type ListWorkflowRequest struct {
	Header       int
	Names        []string
	Descriptions []string
	StartTime    int64
	EndTime      int64
	CurrentPage  uint64
	PageSize     uint64
}

type ListWorkflowResponse struct {
	WorkflowStats []string `json:"WorkflowStats"`
}

func decodeListWorkflowRequest(r *http.Request) (interface{}, error) {
	pathParams := mux.Vars(r)
	_ = pathParams
	queryParams := r.URL.Query()
	_ = queryParams
	return &ListWorkflowRequest{}, nil
}

func (w *workFlowAPI) ListWorkflow(ctx context.Context, req interface{}) (interface{}, error) {
	request, ok := req.(*ListWorkflowRequest)
	if !ok {
		return nil, ErrorBadRequest
	}

	_ = request

	resp := &ListWorkflowResponse{}

	return resp, nil
}

type DeleteWorkflowRequest struct {
	Header       int
	Names        []string
	Descriptions []string
	StartTime    int64
	EndTime      int64
	CurrentPage  uint64
	PageSize     uint64
}

type DeleteWorkflowResponse struct {
	WorkflowStats []string `json:"WorkflowStats"`
}

func decodeDeleteWorkflowRequest(r *http.Request) (interface{}, error) {
	pathParams := mux.Vars(r)
	_ = pathParams
	queryParams := r.URL.Query()
	_ = queryParams
	return &DeleteWorkflowRequest{}, nil
}

func (w *workFlowAPI) DeleteWorkflow(ctx context.Context, req interface{}) (interface{}, error) {
	request, ok := req.(*DeleteWorkflowRequest)
	if !ok {
		return nil, ErrorBadRequest
	}

	_ = request

	resp := &DeleteWorkflowResponse{}

	return resp, nil
}
