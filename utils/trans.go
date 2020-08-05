package utils

func ToString(v interface{}) string {
	if ret, ok := v.(string); ok {
		return ret
	}
	return ""
}

func ToBytes(v interface{}) []byte {
	if ret, ok := v.([]byte); ok {
		return ret
	}
	return nil
}

func ToBool(v interface{}) bool {
	if ret, ok := v.(bool); ok {
		return ret
	}
	return false
}

func ToFloat64(v interface{}) float64 {
	if ret, ok := v.(float64); ok {
		return ret
	}
	return 0
}

func ToInt64(v interface{}) int64 {
	if ret, ok := v.(int64); ok {
		return ret
	}
	return 0
}
