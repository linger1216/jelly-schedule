package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/rpc/json"
	"log"
	"net/http"
)

func checkError(err error) {
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func main() {
	url := "http://localhost:12345/rpc"
	args := "hello jsonrpc"

	message, err := json.EncodeClientRequest("DefaultJsonRPCService.Exec", args)
	checkError(err)

	resp, err := http.Post(url, "application/json", bytes.NewReader(message))
	defer resp.Body.Close()

	checkError(err)

	reply := ""
	err = json.DecodeClientResponse(resp.Body, &reply)
	checkError(err)

	fmt.Printf("%s\n", reply)
}
