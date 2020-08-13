package core

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/linger1216/jelly-schedule/utils"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
	"github.com/linger1216/go-pipeline/pipe"
)

const (
	HeadCountKey = "X-Total-Count"
)

type workFlowAPI struct {
	db                   *sqlx.DB
	createWorkflowFilter pipe.Filter
}

func NewWorkflowAPI(db *sqlx.DB) *workFlowAPI {
	_, err := db.Exec(createWorkflowTableSql())
	if err != nil {
		panic(err)
	}
	return &workFlowAPI{db: db}
}

type CreateWorkflowRequest struct {
	Workflows []*WorkFlow
}

type CreateWorkflowResponse struct {
	Ids []string `json:"ids"`
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
		return nil, ErrBadRequest
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
	Ids []string `json:"ids"`
}

type GetWorkflowResponse struct {
	Workflows []*WorkFlow `json:"workflows"`
}

func decodeGetWorkflowRequest(r *http.Request) (interface{}, error) {
	pathParams := mux.Vars(r)
	_ = pathParams
	queryParams := r.URL.Query()
	_ = queryParams

	req := &GetWorkflowRequest{}
	arr := pathParams["ids"]
	req.Ids = strings.Split(arr, ",")
	return req, nil
}

func (w *workFlowAPI) GetWorkflow(ctx context.Context, req interface{}) (interface{}, error) {
	in, ok := req.(*GetWorkflowRequest)
	if !ok {
		return nil, ErrBadRequest
	}

	resp := &GetWorkflowResponse{}
	resp.Workflows = make([]*WorkFlow, len(in.Ids))
	actualIds := make([]string, 0, len(in.Ids))
	for i := range in.Ids {
		if len(in.Ids[i]) > 0 {
			actualIds = append(actualIds, in.Ids[i])
		} else {
			resp.Workflows[i] = &WorkFlow{}
		}
	}

	query, args := getWorkflowSql(actualIds)
	l.Debug(query)

	workflows, err := w.queryWorkflow(query, args)
	if err != nil {
		return nil, err
	}

	pos := 0
	for i := range resp.Workflows {
		if resp.Workflows[i] == nil {
			if pos < len(workflows) && workflows[pos] != nil {
				resp.Workflows[i] = workflows[pos]
				pos++
			} else {
				resp.Workflows[i] = &WorkFlow{}
			}
		}
	}
	return resp, nil
}

type UpdateWorkflowRequest struct {
	WorkFlows []*WorkFlow `json:"workflows"`
}

type UpdateWorkflowResponse struct {
}

func decodeUpdateWorkflowRequest(r *http.Request) (interface{}, error) {
	req := &UpdateWorkflowRequest{}
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
		if err = jsoniter.ConfigFastest.Unmarshal(buf, &req.WorkFlows); err != nil {
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

func (w *workFlowAPI) UpdateWorkflow(ctx context.Context, req interface{}) (interface{}, error) {
	in, ok := req.(*UpdateWorkflowRequest)
	if !ok || len(in.WorkFlows) == 0 {
		return nil, ErrorInvalidPara
	}
	query, args, err := upsertWorkflowSql(in.WorkFlows)
	if err != nil {
		return nil, err
	}

	l.Debugf(query)
	_, err = w.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return &UpdateWorkflowResponse{}, nil
}

type ListWorkflowRequest struct {
	Header      int      `json:"header"`
	Names       []string `json:"names"`
	States      []string `json:"states"`
	StartTime   int64    `json:"startTime"`
	EndTime     int64    `json:"endTime"`
	CurrentPage uint64   `json:"currentPage"`
	PageSize    uint64   `json:"pageSize"`
}

type ListWorkflowResponse struct {
	Headers   []*KV       `json:"headers,omitempty"`
	Workflows []*WorkFlow `json:"workflows,omitempty"`
}

func decodeListWorkflowRequestHeader(r *http.Request) (interface{}, error) {
	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	req := decodeListQueryParams(nil, queryParams)
	req.Header = 1
	return req, nil
}

func decodeListWorkflowRequest(r *http.Request) (interface{}, error) {
	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	req := decodeListQueryParams(nil, queryParams)
	req.Header = 0
	return req, nil
}

func (w *workFlowAPI) ListWorkflow(ctx context.Context, req interface{}) (interface{}, error) {
	in, ok := req.(*ListWorkflowRequest)
	if !ok {
		return nil, ErrBadRequest
	}

	if in.PageSize == 0 {
		in.PageSize = 100
	}

	resp := &ListWorkflowResponse{}
	query, args := listAssetsSql(in)
	l.Debugf(query)

	if in.Header > 0 {
		count := int64(0)
		err := w.db.Get(&count, query, args...)
		if err != nil {
			return nil, err
		}
		if count == 0 {
			return nil, ErrNotFound
		}
		resp.Headers = append(resp.Headers, &KV{
			Key:   HeadCountKey,
			Value: utils.Int64ToString(count),
		})
	} else {
		ret, err := w.queryWorkflow(query, args)
		if err != nil {
			return nil, err
		}
		resp.Workflows = ret
	}
	return resp, nil
}

func decodeListQueryParams(req *ListWorkflowRequest, queryParams url.Values) *ListWorkflowRequest {
	if req == nil {
		req = &ListWorkflowRequest{}
	}

	if str, ok := queryParams["header"]; ok && len(str) > 0 {
		size, _ := strconv.Atoi(str[0])
		req.Header = size
	}

	if arr, ok := queryParams["names"]; ok {
		req.Names = make([]string, 0, len(arr))
		for i := range arr {
			if len(arr[i]) > 0 {
				req.Names = append(req.Names, arr[i])
			}
		}
	}

	if arr, ok := queryParams["descriptions"]; ok {
		req.States = make([]string, 0, len(arr))
		for i := range arr {
			if len(arr[i]) > 0 {
				req.States = append(req.States, arr[i])
			}
		}
	}

	if str, ok := queryParams["start_time"]; ok && len(str) > 0 {
		req.StartTime = utils.StringToInt64(str[0])
	}
	if str, ok := queryParams["end_time"]; ok && len(str) > 0 {
		req.EndTime = utils.StringToInt64(str[0])
	}

	if str, ok := queryParams["current_page"]; ok && len(str) > 0 {
		page, _ := strconv.Atoi(str[0])
		req.CurrentPage = uint64(page)
	}

	if str, ok := queryParams["page_size"]; ok && len(str) > 0 {
		size, _ := strconv.Atoi(str[0])
		req.PageSize = uint64(size)
	}

	return req
}

type DeleteWorkflowRequest struct {
	Ids []string `json:"ids"`
}

type DeleteWorkflowResponse struct {
}

func decodeDeleteWorkflowRequest(r *http.Request) (interface{}, error) {
	pathParams := mux.Vars(r)
	_ = pathParams
	queryParams := r.URL.Query()
	_ = queryParams

	req := &DeleteWorkflowRequest{}
	arr := pathParams["ids"]
	req.Ids = strings.Split(arr, ",")
	return req, nil
}

func (w *workFlowAPI) DeleteWorkflow(ctx context.Context, req interface{}) (interface{}, error) {
	in, ok := req.(*DeleteWorkflowRequest)
	if !ok {
		return nil, ErrBadRequest
	}

	resp := &DeleteWorkflowResponse{}
	if len(in.Ids) == 0 {
		return resp, nil
	}

	query, args := deleteWorkflowSql(in.Ids)
	l.Debug(query)
	_, err := w.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// private
func (w *workFlowAPI) queryWorkflow(query string, args []interface{}) ([]*WorkFlow, error) {
	rows, err := w.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	ret := make([]*WorkFlow, 0)
	for rows.Next() {
		line := make(map[string]interface{})
		err = rows.MapScan(line)
		if err != nil {
			return nil, err
		}
		if tc, err := transWorkflow("", line); err == nil && tc != nil {
			ret = append(ret, tc)
		}
	}
	if len(ret) == 0 {
		return nil, nil
	}
	return ret, nil
}

func encodeHTTPWorkflowResponse(w http.ResponseWriter, response interface{}) error {
	encoder := jsoniter.ConfigFastest.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	switch x := response.(type) {
	case *CreateWorkflowResponse:
		return encoder.Encode(x.Ids)
	case *GetWorkflowResponse:
		return encoder.Encode(x.Workflows)
	case *ListWorkflowResponse:
		if len(x.Headers) > 0 {
			return writeHeader(w, x.Headers)
		}
		return encoder.Encode(x.Workflows)
	}
	return encoder.Encode(response)
}

func transWorkflow(prefix string, m map[string]interface{}) (*WorkFlow, error) {
	ret := &WorkFlow{}
	if v, ok := m[prefix+"id"]; ok {
		ret.Id = utils.ToString(v)
	}

	if v, ok := m[prefix+"name"]; ok {
		ret.Name = utils.ToString(v)
	}

	if v, ok := m[prefix+"description"]; ok {
		ret.Description = utils.ToString(v)
	}

	if v, ok := m[prefix+"job_ids"]; ok {
		err := jsoniter.ConfigFastest.Unmarshal([]byte(utils.ToString(v)), &ret.JobIds)
		if err != nil {
			return nil, err
		}
	}

	if v, ok := m[prefix+"cron"]; ok {
		ret.Cron = utils.ToString(v)
	}

	if v, ok := m[prefix+"para"]; ok {
		ret.Para = utils.ToString(v)
	}

	if v, ok := m[prefix+"execute_limit"]; ok {
		ret.ExecuteLimit = utils.ToInt64(v)
	}

	if v, ok := m[prefix+"error_policy"]; ok {
		ret.ErrorPolicy = utils.ToString(v)
	}

	if v, ok := m[prefix+"belong_executor"]; ok {
		ret.BelongExecutor = utils.ToString(v)
	}

	if v, ok := m[prefix+"state"]; ok {
		ret.State = utils.ToString(v)
	}

	if v, ok := m[prefix+"create_time"]; ok {
		ret.CreateTime = utils.ToInt64(v)
	}

	if v, ok := m[prefix+"update_time"]; ok {
		ret.UpdateTime = utils.ToInt64(v)
	}
	return ret, nil
}
