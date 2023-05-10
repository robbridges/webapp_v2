package controllers

import "net/http"

type Template interface {
	Execute(w http.ResponseWriter, r *http.Request, data interface{})
}

type MockTemplate struct {
	ExecuteFunc func(w http.ResponseWriter, r *http.Request, data interface{})
}

func (mt *MockTemplate) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	mt.ExecuteFunc(w, r, data)
}
