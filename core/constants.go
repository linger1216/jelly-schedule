package core

import (
	"errors"
	"net/http"
)

const (
	StateAvaiable  = "avaiable"
	StateExecuting = "executing"
	StateFinish    = "finish"

	ErrPolicyPanic  = "panic"
	ErrPolicyIgnore = "ignore"
	ErrPolicyRetry  = "retry"

	ExecUnlimitCount = -1
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
