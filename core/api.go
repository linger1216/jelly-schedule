package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/cors"
)

type StatusCoder interface {
	StatusCode() int
}

type Header interface {
	Headers() http.Header
}

type apiError struct {
	code int
	msg  string
}

func newApiError(code int, msg string) *apiError {
	return &apiError{code: code, msg: msg}
}

func (a *apiError) Error() string {
	return a.msg
}

func (a *apiError) RuntimeError() {
	panic("implement me")
}

func (a *apiError) StatusCode() int {
	return a.code
}

type HandleFunc func(w http.ResponseWriter, r *http.Request)
type DecodeRequestFunc func(r *http.Request) (interface{}, error)
type encodeResponseFunc func(w http.ResponseWriter, response interface{}) error

type KV struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func writeHeader(w http.ResponseWriter, kvs []*KV) error {
	for i := range kvs {
		w.Header().Set(kvs[i].Key, kvs[i].Value)
	}
	return nil
}

func encodeHTTPGenericResponse(w http.ResponseWriter, response interface{}) error {
	encoder := jsoniter.ConfigFastest.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(response)
}

func DefaultErrorEncoder(err error, w http.ResponseWriter) {
	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
	if marshaler, ok := err.(json.Marshaler); ok {
		if jsonBody, marshalErr := marshaler.MarshalJSON(); marshalErr == nil {
			contentType, body = "application/json; charset=utf-8", jsonBody
		}
	}
	w.Header().Set("Content-Type", contentType)
	if headerer, ok := err.(Header); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := http.StatusInternalServerError
	if sc, ok := err.(StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	w.Write(body)
}

func HandleFuncWrapper(dec DecodeRequestFunc, e Endpoint, enc encodeResponseFunc) HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := dec(r)
		if err != nil {
			DefaultErrorEncoder(err, w)
			return
		}
		resp, err := e(context.Background(), req)
		if err != nil {
			DefaultErrorEncoder(err, w)
			return
		}
		enc(w, resp)
	}
}

type scheduleAPI struct {
	job     *jobAPI
	worflow *workFlowAPI
}

func NewScheduleAPI(etcd *Etcd, db *sqlx.DB) *scheduleAPI {
	api := &scheduleAPI{job: NewJobAPI(etcd), worflow: NewWorkflowAPI(db)}
	return api
}

func (w *scheduleAPI) Start(port int) error {
	m := mux.NewRouter()

	// 获得Job List
	m.HandleFunc("/schedule/job",
		HandleFuncWrapper(decodeGetJobListRequest, w.job.getJobList, encodeHTTPGenericResponse)).
		Methods("GET")

	// 获得Job
	m.HandleFunc("/schedule/job/{ids}",
		HandleFuncWrapper(decodeGetJobRequest, w.job.getJob, encodeHTTPGenericResponse)).
		Methods("GET")

	// 创建工作流
	m.HandleFunc("/schedule/workflow",
		HandleFuncWrapper(decodeCreateWorkflowRequest, w.worflow.CreateWorkflow, encodeHTTPWorkflowResponse)).
		Methods("POST")

	// 删除工作流
	m.HandleFunc("/schedule/workflow",
		HandleFuncWrapper(decodeDeleteWorkflowRequest, w.worflow.DeleteWorkflow, encodeHTTPWorkflowResponse)).
		Methods("DELETE")

	// 更改工作流
	m.HandleFunc("/schedule/workflow",
		HandleFuncWrapper(decodeUpdateWorkflowRequest, w.worflow.UpdateWorkflow, encodeHTTPWorkflowResponse)).
		Methods("PUT")

	// List 工作流
	m.HandleFunc("/schedule/workflow",
		HandleFuncWrapper(decodeListWorkflowRequest, w.worflow.ListWorkflow, encodeHTTPWorkflowResponse)).
		Methods("GET")

	// HEAD 工作流
	m.HandleFunc("/schedule/workflow",
		HandleFuncWrapper(decodeListWorkflowRequestHeader, w.worflow.ListWorkflow, encodeHTTPWorkflowResponse)).
		Methods("HEAD")

	// get 工作流
	m.HandleFunc("/schedule/workflow/{ids}",
		HandleFuncWrapper(decodeGetWorkflowRequest, w.worflow.GetWorkflow, encodeHTTPWorkflowResponse)).
		Methods("GET")

	return http.ListenAndServe(fmt.Sprintf(":%d", port), cors.AllowAll().Handler(m))
}
