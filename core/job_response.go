package core

import (
	"bytes"
	"strings"
)

const (
	Separate = ";"
)

type MergeFunc func(paras ...string) string
type SplitFunc func(paras string) []string

func MergeFactory(sep string) MergeFunc {
	return func(paras ...string) string {
		return _merge(sep, paras...)
	}
}

func SplitFactory(sep string) SplitFunc {
	return func(paras string) []string {
		return _split(sep, paras)
	}
}

func _split(sep string, paras string) []string {
	arr := strings.Split(paras, sep)
	ret := make([]string, len(arr))
	for i := range arr {
		ret[i] = arr[i]
	}
	return ret
}

func _merge(sep string, paras ...string) string {
	var buf bytes.Buffer
	for i := range paras {
		buf.WriteString(paras[i])
		if i < len(paras)-1 {
			buf.WriteString(sep)
		}
	}
	return buf.String()
}

//
//func exactSerialRequest(req string) string {
//	var arg string
//	switch x := req.(type) {
//	case string:
//		arg = x
//	case []string:
//		if len(x) == 1 {
//			arg = x[0]
//		} else {
//			arg = x
//		}
//	}
//	return arg
//}

//
//func exactParallelRequest(req string, size int) ([]string, error) {
//	args := make([]string, 0, size)
//	switch x := req.(type) {
//	case string:
//		arr := make([]string, 0)
//		_ = jsoniter.ConfigFastest.Unmarshal([]byte(x), &arr)
//		if len(arr) == 0 {
//			arr = strings.Split(x, Separate)
//		}
//		if len(arr) != size {
//			return nil, ErrorJobParaInvalid
//		}
//		for i := range arr {
//			args = append(args, arr[i])
//		}
//	case []string:
//		// 后面是一个单独的任务, 要保证后续的参数正常
//		if size > 1 && len(x) != size {
//			return nil, ErrorJobParaInvalid
//		} else {
//			args = x
//		}
//	}
//	return args, nil
//}
//

//
//func ExactJobRequests(req string) ([]string, error) {
//	args := make([]string, 0)
//	switch x := req.(type) {
//	case string:
//		_ = jsoniter.ConfigFastest.Unmarshal([]byte(x), &args)
//		if len(args) > 0 {
//			break
//		}
//		arr := strings.Split(x, Separate)
//		for i := range arr {
//			args = append(args, arr[i])
//		}
//	case []string:
//		for i := range x {
//			if v, ok := x[i].(string); ok {
//				args = append(args, v)
//			} else {
//				return nil, fmt.Errorf("%d arg %v invalid", i, x[i])
//			}
//		}
//	}
//	return args, nil
//}
