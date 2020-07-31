package core

//
//import (
//	"context"
//	"github.com/gorilla/mux"
//	jsoniter "github.com/json-iterator/go"
//	"net/http"
//)
//
//type getWorkerListRequest struct{}
//
//type getWorkerListResponse struct {
//	WorkerStats []*JobStats `json:"workStats"`
//}
//
//func decodeGetWorkerListRequest(r *http.Request) (interface{}, error) {
//	pathParams := mux.Vars(r)
//	_ = pathParams
//	queryParams := r.URL.Query()
//	_ = queryParams
//	return &getWorkerListRequest{}, nil
//}
//
//func (w *scheduleAPI) getWorkerList(ctx context.Context, req interface{}) (interface{}, error) {
//	request, ok := req.(*getWorkerListRequest)
//	if !ok {
//		return nil, ErrorBadRequest
//	}
//
//	_ = request
//
//	_, v, err := w.etcd.GetWithPrefixKey(ctx, WorkerPrefix)
//	if err != nil {
//		return nil, err
//	}
//
//	resp := &getWorkerListResponse{}
//	for i := range v {
//		stats := &WorkerStats{}
//		err = jsoniter.ConfigFastest.Unmarshal([]byte(v[i]), stats)
//		if err != nil {
//			return nil, err
//		}
//		resp.WorkerStats = append(resp.WorkerStats, stats)
//	}
//
//	return resp, nil
//}
//
//type Resp struct {
//	ID   string
//	Name string
//}
//
//func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
//
//	resp := &Resp{
//		ID:   "ghhghg",
//		Name: "kknjn",
//	}
//
//	_ = encodeHTTPGenericResponse(w, resp)
//}
