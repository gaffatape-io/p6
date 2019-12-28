package rest

import (
	"cloud.google.com/go/firestore"
	"context"
	errs "github.com/gaffatape-io/gopherrs"
	"github.com/gaffatape-io/p6/crud"
	"k8s.io/klog"
	"net/http"
	"strings"
)

func registerObjectiveHandlers(s *crud.Store, mux *http.ServeMux) {
	oh := objectiveHttpHandler(s, s.RunTx)
	mux.HandleFunc("/o/", oh)
	mux.HandleFunc("/o", oh)
}

func objectiveHttpHandler(s crud.ObjectiveStore, runTx crud.TxRun) http.HandlerFunc {
	h := &objectiveHandler{s, runTx}

	return func(w http.ResponseWriter, r *http.Request) {
		klog.V(11).Info(r.Method, " ", r.URL.Path)
		handler := func(r *http.Request) (interface{}, error) {
			klog.Error("Not allowed ", r.Method, r.URL.Path)
			return nil, errs.InvalidArgumentf(nil, "%s not supported", r.Method)
		}

		switch {
		case r.Method == http.MethodPut:
			handler = h.put

		case r.Method == http.MethodPost:
			handler = h.post
		}

		resp, err := handler(r)
		if err != nil {
			writeStatus(w, err)
			return
		}

		err = writeJson(w, resp)
		if err != nil {
			klog.Error("Failed to write response; result unsure")
		}
	}
}

type Objective struct {
	HItem
}

type objectiveHandler struct {
	s     crud.ObjectiveStore
	runTx crud.TxRun
}

func (h *objectiveHandler) createObjectiveEntity(ctx context.Context, o Objective) (crud.ObjectiveEntity, error) {
	data := crud.Objective{crud.HItem{crud.Item{o.Summary, o.Description}, o.ParentID}}

	if data.ParentID == "" {
		return h.s.CreateObjective(ctx, nil, data)
	}

	var entity crud.ObjectiveEntity
	err := h.runTx(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		_, err := h.s.ReadObjective(ctx, tx, data.ParentID)
		if err != nil {
			return err
		}

		entity, err = h.s.CreateObjective(ctx, tx, data)
		return err
	})

	return entity, err
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

	if o.Summary == "" {
		return nil, errs.InvalidArgumentf(nil, "summary not set")
	}

	ctx := r.Context()
	entity, err := h.createObjectiveEntity(ctx, o)
	if errs.IsNotFound(err) {
		return nil, errs.InvalidArgumentf(nil, "parent not found")
	} else if err != nil {
		return nil, errs.Internal(err)
	}

	klog.Info("objective created:", entity.ID)
	return &Objective{HItem{Item{entity.ID, entity.Summary, entity.Description}, entity.ParentID}}, nil
}

func (h *objectiveHandler) post(r *http.Request) (interface{}, error) {
	var o Objective
	err := readJson(r.Body, &o)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
