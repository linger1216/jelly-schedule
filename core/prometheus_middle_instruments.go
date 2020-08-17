package core

//
//import (
//	"context"
//	"github.com/go-kit/kit/endpoint"
//	"time"
//)
//
//
//// InstrumentingMiddleware returns an endpoint middleware that records
//// the duration of each invocation to the passed histogram. The middleware adds
//// a single field: "success", which is "true" if no error is returned, and
//// "false" otherwise.
//
//func Instrumenting(duration Histogram, counter Counter) Middleware {
//	return func(next Endpoint) Endpoint {
//		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//			defer func(begin time.Time) {
//				//errString := getErrorString(err)
//				//accessKey := getAccessKey(ctx)
//				//count := getCount(request)
//				counter.With("error", errString, "access_key", accessKey).Add(count)
//				duration.With("error", errString, "access_key", accessKey).Observe(time.Since(begin).Seconds())
//			}(time.Now())
//			return next(ctx, request)
//		}
//	}
//}

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
