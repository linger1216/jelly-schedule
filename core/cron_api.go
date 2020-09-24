package core

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/linger1216/jelly-schedule/utils"
	"github.com/robfig/cron/v3"
	"net/http"
	"time"
)

type cronAPI struct {
}

func NewCronAPI() *cronAPI {
	return &cronAPI{}
}

type GetCronRequest struct {
	Expr      string `json:"expr,omitempty" yaml:"expr" `
	NextCount int    `json:"next_count,omitempty" yaml:"next_count" `
}

type GetCronResponse struct {
	Expr  string   `json:"expr,omitempty" yaml:"expr" `
	Nexts []string `json:"nexts,omitempty" yaml:"nexts" `
}

func decodeGetCronRequest(r *http.Request) (interface{}, error) {
	pathParams := mux.Vars(r)
	_ = pathParams
	queryParams := r.URL.Query()
	_ = queryParams

	req := &GetCronRequest{}
	if arr, ok := queryParams["expr"]; ok && len(arr) > 0 {
		req.Expr = arr[0]
	}
	if str, ok := queryParams["next_count"]; ok && len(str) > 0 {
		req.NextCount = int(utils.StringToInt64(str[0]))
	}
	return req, nil
}

func (w *cronAPI) GetCron(ctx context.Context, req interface{}) (interface{}, error) {
	in, ok := req.(*GetCronRequest)
	if !ok {
		return nil, ErrBadRequest
	}

	sched, err := cron.ParseStandard(in.Expr)
	if err != nil {
		return nil, ErrBadCronExpr
	}

	resp := &GetCronResponse{}
	resp.Expr = in.Expr

	t := time.Now().Add(-1 * time.Second)
	for i := 0; i < in.NextCount; i++ {
		t = sched.Next(t)
		resp.Nexts = append(resp.Nexts, t.Format("2006-01-02 15:04:05"))
	}
	return resp, nil
}
