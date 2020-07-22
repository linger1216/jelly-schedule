package jsonrpc

import (
	"fmt"
	"net/http"
)

type Request struct {
	Req string
}

type Response struct {
	Resp string
}

type DefaultJsonRPCService struct{}

func (t *DefaultJsonRPCService) Exec(r *http.Request, arg *Request, result *Response) error {
	//log.Printf("Multiply %d with %d\n", args.A, args.B)
	//*result = args.A * args.B
	fmt.Printf("Req:%s\n", arg.Req)
	result.Resp = fmt.Sprintf("Resp:%s", arg.Req)
	return nil
}
