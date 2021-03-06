package core

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
