package core

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/gorilla/mux"
//	jsoniter "github.com/json-iterator/go"
//	"github.com/rs/cors"
//	"net/http"
//)
//
//type HandleFunc func(w http.ResponseWriter, r *http.Request)
//type DecodeRequestFunc func(r *http.Request) (interface{}, error)
//type encodeResponseFunc func(w http.ResponseWriter, response interface{}) error
//
//func encodeHTTPGenericResponse(w http.ResponseWriter, response interface{}) error {
//	encoder := jsoniter.ConfigFastest.NewEncoder(w)
//	encoder.SetEscapeHTML(false)
//	return encoder.Encode(response)
//}
//
//func DefaultErrorEncoder(err error, w http.ResponseWriter) {
//	contentType, body := "text/plain; charset=utf-8", []byte(err.Error())
//	if marshaler, ok := err.(json.Marshaler); ok {
//		if jsonBody, marshalErr := marshaler.MarshalJSON(); marshalErr == nil {
//			contentType, body = "application/json; charset=utf-8", jsonBody
//		}
//	}
//	w.Header().Set("Content-Type", contentType)
//	if headerer, ok := err.(Header); ok {
//		for k, values := range headerer.Headers() {
//			for _, v := range values {
//				w.Header().Add(k, v)
//			}
//		}
//	}
//	code := http.StatusInternalServerError
//	if sc, ok := err.(StatusCoder); ok {
//		code = sc.StatusCode()
//	}
//	w.WriteHeader(code)
//	w.Write(body)
//}
//
//func HandleFuncWrapper(dec DecodeRequestFunc, e Endpoint, enc encodeResponseFunc) HandleFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		req, err := dec(r)
//		if err != nil {
//			DefaultErrorEncoder(err, w)
//			return
//		}
//		resp, err := e(context.Background(), req)
//		if err != nil {
//			DefaultErrorEncoder(err, w)
//			return
//		}
//		enc(w, resp)
//	}
//}
//
//type scheduleAPI struct {
//	etcd *Etcd
//}
//
//func NewScheduleAPI(etcd *Etcd) *scheduleAPI {
//	return &scheduleAPI{etcd: etcd}
//}
//
//func (w *scheduleAPI) Start(port int) error {
//	m := mux.NewRouter()
//
//	// 获得work节点
//	m.HandleFunc("/schedule/worker",
//		HandleFuncWrapper(decodeGetWorkerListRequest, w.getWorkerList, encodeHTTPGenericResponse)).
//		Methods("GET")
//
//	// 获得leader
//	m.HandleFunc("/schedule/leader", ArticlesCategoryHandler).Methods("GET")
//
//	// 获得所有job
//	m.HandleFunc("/schedule/job", ArticlesCategoryHandler).Methods("GET")
//
//	// 获得job
//	m.HandleFunc("/schedule/job/{id}", ArticlesCategoryHandler).Methods("GET")
//
//	// 删除job
//	m.HandleFunc("/schedule/job/{id}", ArticlesCategoryHandler).Methods("DELETE")
//
//	// 获得所有workflow
//	m.HandleFunc("/schedule/workflow", ArticlesCategoryHandler).Methods("GET")
//
//	// 获得workflow
//	m.HandleFunc("/schedule/workflow/{id}", ArticlesCategoryHandler).Methods("GET")
//
//	// 删除workflow
//	m.HandleFunc("/schedule/workflow/{id}", ArticlesCategoryHandler).Methods("DELETE")
//
//	return http.ListenAndServe(fmt.Sprintf(":%d", port), cors.AllowAll().Handler(m))
//}
