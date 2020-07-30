package core

//
//import (
//	"context"
//	"github.com/gorilla/mux"
//	jsoniter "github.com/json-iterator/go"
//	"github.com/rs/cors"
//	"net/http"
//)
//
//// todo
//// need config
//const (
//	DefaultRestfulPort = ":35744"
//)
//
//type scheduleAPI struct {
//	etcd *etcdv3.Etcd
//}
//
//func (w *scheduleAPI) Start() error {
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
//	return http.ListenAndServe(DefaultRestfulPort, cors.AllowAll().Handler(m))
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
//	_, v, err := w.etcd.GetWithPrefixKey(WorkerPrefix)
//	if err != nil {
//		return nil, err
//	}
//
//	resp := &getWorkerListResponse{}
//	for i := range v {
//		stats := &WorkerStats{}
//		err = jsoniter.ConfigFastest.Unmarshal(v[i], stats)
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
