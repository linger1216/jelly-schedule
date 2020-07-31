package core

import "net/http"

type StatusCoder interface {
	StatusCode() int
}

type Header interface {
	Headers() http.Header
}

type apiError struct {
	code int
	msg  string
}

func newApiError(code int, msg string) *apiError {
	return &apiError{code: code, msg: msg}
}

func (a *apiError) Error() string {
	return a.msg
}

func (a *apiError) RuntimeError() {
	panic("implement me")
}

func (a *apiError) StatusCode() int {
	return a.code
}
