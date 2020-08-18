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
	config  *HttpConfig
}

type HttpConfig struct {
	Port      int `json:"port" yaml:"port" `
	PProfPort int `json:"pprofPort" yaml:"pprofPort" `
}

func NewScheduleAPI(etcd *Etcd, db *sqlx.DB, config *HttpConfig) *scheduleAPI {
	api := &scheduleAPI{job: NewJobAPI(etcd), worflow: NewWorkflowAPI(db), config: config}
	return api
}

func (w *scheduleAPI) Start() error {

	if w.config.PProfPort > 0 {
		go func() {
			l.Debugf("pprof start: %d", w.config.PProfPort)
			http.ListenAndServe(fmt.Sprintf(":%d", w.config.PProfPort), nil)
		}()
	}

	m := mux.NewRouter()

	// 获得Job List
	m.HandleFunc("/schedule/job",
		HandleFuncWrapper(decodeGetJobListRequest, w.job.listJob, encodeHTTPJobResponse)).
		Methods("GET")

	// 获得Job
	m.HandleFunc("/schedule/job/{ids}",
		HandleFuncWrapper(decodeGetJobRequest, w.job.getJob, encodeHTTPJobResponse)).
		Methods("GET")

	// create工作流
	m.HandleFunc("/schedule/workflow",
		HandleFuncWrapper(decodeCreateWorkflowRequest, w.worflow.CreateWorkflow, encodeHTTPWorkflowResponse)).
		Methods("POST")

	// del工作流
	m.HandleFunc("/schedule/workflow",
		HandleFuncWrapper(decodeDeleteWorkflowRequest, w.worflow.DeleteWorkflow, encodeHTTPWorkflowResponse)).
		Methods("DELETE")

	// update工作流
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

	l.Debugf("api start: %d", w.config.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", w.config.Port), cors.AllowAll().Handler(m))
}
