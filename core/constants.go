package core

import (
	"context"
	"errors"
	"net/http"
)

type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)
type Middleware func(Endpoint) Endpoint

const (
	StateAvaiable  = "avaiable"
	StateExecuting = "executing"
	StateFailed    = "failed"
	StateFinish    = "finish"

	ErrPolicyReturn = "return"
	ErrPolicyIgnore = "ignore"
	ErrPolicyRetry  = "retry"

	ExecUnlimitCount  = -1
	DefaultRetryCount = 3
)

var (
	// etcd
	ErrKeyAlreadyExists  = errors.New("key already exists")
	ErrEtcdLeaseNotFound = errors.New("lease not found")

	// api
	ErrBadRequest    = newApiError(http.StatusBadRequest, "StatusBadRequest")
	ErrorInvalidPara = newApiError(http.StatusBadRequest, "ErrorInvalidPara")
	ErrNotFound      = newApiError(http.StatusNotFound, "ErrNotFound")
)
