package core

import "errors"

var (
	ErrKeyAlreadyExists  = errors.New("key already exists")
	ErrEtcdLeaseNotFound = errors.New("lease not found")
	//ErrInsertKV  = errors.New("insert kv error")
)
