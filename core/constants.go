package core

import (
	"errors"
	"net/http"
)

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

	// prometheus
	PrometheusNamespace = "Jelly"
	PrometheusSubsystem = "Schedule"

	RemoteServerMethod = `JsonRPCService.Exec`
)

var (
	// etcd
	ErrKeyAlreadyExists  = errors.New("key already exists")
	ErrEtcdLeaseNotFound = errors.New("lease not found")
	ErrJobNotFound       = errors.New("job not found")
	ErrorJobParaInvalid  = errors.New("job para invalid")

	// api
	ErrBadRequest    = newApiError(http.StatusBadRequest, "StatusBadRequest")
	ErrorInvalidPara = newApiError(http.StatusBadRequest, "ErrorInvalidPara")
	ErrNotFound      = newApiError(http.StatusNotFound, "ErrNotFound")
)
