package rest

import (
	errs "github.com/gaffatape-io/gopherrs"
	"github.com/gaffatape-io/p6/crud"
	"github.com/gaffatape-io/p6/okrs"
	"k8s.io/klog"
	"net/http"
	"strings"
)

func registerObjectiveHandlers(o *okrs.Objectives, mux *http.ServeMux) {
	oh := objectiveHttpHandler(o)
	mux.HandleFunc("/o/", oh)
	mux.HandleFunc("/o", oh)
}

func objectiveHttpHandler(o *okrs.Objectives) http.HandlerFunc {
	h := &objectiveHandler{o}
	d := &methodDispatcher{
		put:  h.put,
		post: h.post,
	}

	return d.dispatch
}

type Objective crud.Objective
type ObjectiveEntity crud.ObjectiveEntity

type objectiveHandler struct {
	o *okrs.Objectives
}

func (h *objectiveHandler) put(r *http.Request) (interface{}, error) {
	if !strings.HasSuffix(r.URL.Path, "/o") {
		return nil, errs.InvalidArgumentf(nil, "invalid path")
	}

	var o Objective
	err := readJson(r.Body, &o)
	if err != nil {
		return nil, errs.InvalidArgumentf(nil, "deserialization failed")
	}

	ctx := r.Context()
	entity, err := h.o.Create(ctx, crud.Objective(o))
	if errs.IsFailedPrecondition(err) {
		return nil, err
	} else if errs.IsInvalidArgument(err) {
		return nil, err
	} else if err != nil {
		return nil, errs.Internal(err)
	}

	klog.Info("objective created:", entity.ID)
	return ObjectiveEntity(entity), nil
}

func (h *objectiveHandler) post(r *http.Request) (interface{}, error) {
	var o Objective
	err := readJson(r.Body, &o)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
