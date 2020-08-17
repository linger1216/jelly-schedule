package core

import (
	"context"
	"time"
)

//
//import (
//	"context"
//	"github.com/go-kit/kit/endpoint"
//	"time"
//)

// InstrumentingMiddleware returns an endpoint middleware that records
// the duration of each invocation to the passed histogram. The middleware adds
// a single field: "success", which is "true" if no error is returned, and
// "false" otherwise.
// progress IGauge

func Instrumenting(latency IHistogram, success, failed ICounter) Middleware {
	return func(next Endpoint) Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			defer func(begin time.Time) {
				latency.Observe(time.Since(begin).Seconds())
			}(time.Now())
			resp, err := next(ctx, request)
			if err != nil {
				failed.Add(1)
			} else {
				success.Add(1)
			}
			return resp, err
		}
	}
}

//
//func getAccessKey(ctx context.Context) string {
//	var ret string
//	if accessKey, ok := ctx.Value("access_key").(string); ok {
//		ret = accessKey
//	}
//	return ret
//}
//
//func getCount(request interface{}) float64 {
//	if r, ok := request.(counter.BatchRequest); ok {
//		return float64(r.GetBatchCount())
//	}
//	return 1
//}
//
//func getErrorString(err error) string {
//	var ret string
//	if err != nil {
//		if e, ok := err.(*core.Error); ok {
//			if e.StatusCode() == 200 {
//				ret = "ok"
//			} else if e.StatusCode() >= 400 && e.StatusCode() < 500 {
//				ret = "client_error"
//			} else {
//				ret = "server_error"
//			}
//		} else {
//			ret = "server_error"
//		}
//	} else {
//		ret = "ok"
//	}
//	return ret
//}
