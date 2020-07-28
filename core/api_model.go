package core

//
//import (
//	"context"
//	"encoding/json"
//	"github.com/gorilla/mux"
//	jsoniter "github.com/json-iterator/go"
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
//type getWorkerListRequest struct{}
//
//type getWorkerListResponse struct {
//	WorkerStats []*WorkerStats `json:"workStats"`
//}
//
//func decodeGetWorkerListRequest(r *http.Request) (interface{}, error) {
//	pathParams := mux.Vars(r)
//	_ = pathParams
//	queryParams := r.URL.Query()
//	_ = queryParams
//	return &getWorkerListRequest{}, nil
//}
