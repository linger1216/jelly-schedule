package core

import (
	"errors"
	"net/http"
)

var (
	// etcd
	ErrKeyAlreadyExists  = errors.New("key already exists")
	ErrEtcdLeaseNotFound = errors.New("lease not found")

	// api
	ErrorBadRequest  = newApiError(http.StatusBadRequest, "StatusBadRequest")
	ErrorInvalidPara = newApiError(http.StatusBadRequest, "ErrorInvalidPara")
)
