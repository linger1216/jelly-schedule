package core

import (
	"errors"
	"net/http"
)

var (
	// etcd
	ErrKeyAlreadyExists  = errors.New("key already exists")
	ErrEtcdLeaseNotFound = errors.New("lease not found")
	//ErrInsertKV  = errors.New("insert kv error")

	// api
	ErrorBadRequest = NewApiError(http.StatusBadRequest, "StatusBadRequest")
)
