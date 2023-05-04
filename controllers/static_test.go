package controllers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockTemplate struct {
	ExecuteFunc func(w http.ResponseWriter, r *http.Request, data interface{})
}

func (mt *MockTemplate) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	mt.ExecuteFunc(w, r, data)
}

func TestStaticHandler(t *testing.T) {
	mockTmpl := &MockTemplate{
		ExecuteFunc: func(w http.ResponseWriter, r *http.Request, data interface{}) {
			w.Write([]byte("Test"))
		},
	}

	handlerFunc := StaticHandler(mockTmpl)
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handlerFunc.ServeHTTP(recorder, req)

	if recorder.Body.String() != "Test" {
		t.Errorf("Expected %q, but got %q", "Test", recorder.Body.String())
	}
}

func TestFAQ(t *testing.T) {
	t.Run("should execute template with questions", func(t *testing.T) {
		mockTpl := &MockTemplate{
			ExecuteFunc: func(w http.ResponseWriter, r *http.Request, data interface{}) {
				questions, ok := data.([]struct {
					Question string
					Answer   template.HTML
				})
				if !ok {
					t.Errorf("expected data to be []struct, but got %T", data)
				}

				if len(questions) != 4 {
					t.Errorf("expected 4 questions, but got %d", len(questions))
				}
			},
		}

		handler := FAQ(mockTpl)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		handler.ServeHTTP(rr, req)
	})
}
