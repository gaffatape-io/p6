package rest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
)

func TestDispatch(t *testing.T) {
	makeHandlerFunc := func(id string) handlerFunc {
		return func(r *http.Request) (interface{}, error) {
			return id, nil
		}
	}

	methods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	dispatcher := &methodDispatcher{
		get:     makeHandlerFunc(http.MethodGet),
		head:    makeHandlerFunc(http.MethodHead),
		post:    makeHandlerFunc(http.MethodPost),
		put:     makeHandlerFunc(http.MethodPut),
		patch:   makeHandlerFunc(http.MethodPatch),
		delete:  makeHandlerFunc(http.MethodDelete),
		connect: makeHandlerFunc(http.MethodConnect),
		options: makeHandlerFunc(http.MethodOptions),
		trace:   makeHandlerFunc(http.MethodTrace),
	}

	for _, m := range methods {
		r := httptest.NewRequest(m, "/", nil)
		w := httptest.NewRecorder()
		dispatcher.dispatch(w, r)
		t.Log(w)

		if w.Code != http.StatusOK {
			t.Fatal()
		}

		body := w.Body.String()
		t.Log(body)
		if !strings.Contains(body, m) {
			t.Fatal()
		}
	}
}

func TestDispatchUnknownMethod(t *testing.T) {
	d := &methodDispatcher{}
	r := httptest.NewRequest("unknown-method", "/", nil)
	w := httptest.NewRecorder()
	d.dispatch(w, r)
	t.Log(w)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatal()
	}
}

func TestDispatch2handlerFunc(t *testing.T) {
	tests := []struct {
		handler handlerFunc
		code    int
	}{
		{
			nil,
			http.StatusMethodNotAllowed,
		},
		{
			func(r *http.Request) (interface{}, error) {
				return "result", nil
			},
			http.StatusOK,
		},
		{
			func(r *http.Request) (interface{}, error) {
				return "ignored", fmt.Errorf("failed")
			},
			http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Log(tc)
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		dispatch2handlerFunc(tc.handler, w, r)
		t.Log(w)

		if w.Code != tc.code {
			t.Fatal()
		}
	}
}
