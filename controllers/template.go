package controllers

import "net/http"

type Template interface {
	Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error)
}

type MockTemplate struct {
	ExecuteFunc func(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error)
}

func (mt *MockTemplate) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	mt.ExecuteFunc(w, r, data)
}
