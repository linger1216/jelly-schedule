package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/linger1216/jelly-schedule/example/jsonrpc"
	"log"
	"net/http"
)

func main() {
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")
	server.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	s := &jsonrpc.DefaultJsonRPCService{}
	err := server.RegisterService(s, "DefaultJsonRPCService")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.Handle("/rpc", server)
	log.Println("JSON RPC service listen and serving on port 12345")
	if err := http.ListenAndServe(":12345", r); err != nil {
		log.Fatalf("Error serving: %s", err)
	}
}
