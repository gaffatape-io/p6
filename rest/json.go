package rest

import (
	"encoding/json"
	"io"
	"k8s.io/klog"
	"net/http"
)

func readJson(r io.Reader, data interface{}) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(&data)
	return err
}

func writeJson(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(data)
}

// handlerFuns defines a request/response style interface for REST handlers.
type handlerFunc func(r *http.Request) (interface{}, error)

// methodDispatcher is a statically types dispatcher for http requests that
// switches from normal http.HandlerFunc to rest.handlerFunc signatures.
type methodDispatcher struct {
	get     handlerFunc
	head    handlerFunc
	post    handlerFunc
	put     handlerFunc
	patch   handlerFunc
	delete  handlerFunc
	connect handlerFunc
	options handlerFunc
	trace   handlerFunc
}

func dispatch2handlerFunc(hf handlerFunc, w http.ResponseWriter, r *http.Request) {
	if hf == nil {
		klog.Error("Not allowed ", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	resp, err := hf(r)
	if err != nil {
		writeStatus(w, err)
		return
	}

	err = writeJson(w, resp)
	if err != nil {
		klog.Error("Failed to write response; result unsure")
	}
}

func (m *methodDispatcher) dispatch(w http.ResponseWriter, r *http.Request) {
	klog.Info(r.Method, " ", r.URL.Path)
	switch {
	case r.Method == http.MethodGet:
		dispatch2handlerFunc(m.get, w, r)

	case r.Method == http.MethodHead:
		dispatch2handlerFunc(m.head, w, r)

	case r.Method == http.MethodPost:
		dispatch2handlerFunc(m.post, w, r)

	case r.Method == http.MethodPut:
		dispatch2handlerFunc(m.put, w, r)

	case r.Method == http.MethodPatch:
		dispatch2handlerFunc(m.patch, w, r)

	case r.Method == http.MethodDelete:
		dispatch2handlerFunc(m.delete, w, r)

	case r.Method == http.MethodConnect:
		dispatch2handlerFunc(m.connect, w, r)

	case r.Method == http.MethodOptions:
		dispatch2handlerFunc(m.options, w, r)

	case r.Method == http.MethodTrace:
		dispatch2handlerFunc(m.trace, w, r)

	default:
		dispatch2handlerFunc(nil, w, r)
	}
}
