package core

import (
	"errors"
	"net/http"
)

const (
	StateAvaiable = "available"
	// available
	StateExecuting = "executing"
	StateFailed    = "failed"
	StateFinish    = "finish"

	// prometheus
	PrometheusNamespace = "Jelly"
	PrometheusSubsystem = "Schedule"
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
