package core

import (
	"testing"
)

func Test_exactParallelRequest(t *testing.T) {
	para := `["\"{\"BeginTimestamp\":1596211200000,\"EndTimestamp\":1596297600000,\"Tokens\":[\"35b277c\"]}\"","\"{\"BeginTimestamp\":1596211200000,\"EndTimestamp\":1596297600000,\"Tokens\":[\"35b277c\"]}\""]`
	exactParallelRequest(para, 2)
}
