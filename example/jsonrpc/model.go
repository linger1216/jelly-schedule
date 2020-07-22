package jsonrpc

import (
	"fmt"
	"net/http"
)

//type Request struct {
//	Req interface{}
//}
//
//type Response struct {
//	Resp interface{}
//}

type Request interface{}
type Response interface{}

type DefaultJsonRPCService struct{}

func (t *DefaultJsonRPCService) Exec(r *http.Request, arg *Request, result *Response) error {
	//log.Printf("Multiply %d with %d\n", args.A, args.B)
	//*result = args.A * args.B
	fmt.Printf("Req:%s\n", (*arg).(string))
	*result = fmt.Sprintf("Resp:%s", (*arg).(string))
	return nil
}
