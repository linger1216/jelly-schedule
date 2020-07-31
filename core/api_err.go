package core

import "net/http"

type StatusCoder interface {
	StatusCode() int
}

type Header interface {
	Headers() http.Header
}

type ApiError struct {
	code int
	msg  string
}

func NewApiError(code int, msg string) *ApiError {
	return &ApiError{code: code, msg: msg}
}

func (a *ApiError) Error() string {
	return a.msg
}

func (a *ApiError) RuntimeError() {
	panic("implement me")
}

func (a *ApiError) StatusCode() int {
	return a.code
}
