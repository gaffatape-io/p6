package rest

import (
	"cloud.google.com/go/firestore"
	"context"
	errs	"github.com/gaffatape-io/gopherrs"
	"github.com/gaffatape-io/p6/crud"
	"k8s.io/klog"
	"net/http"
	"strings"
)

func registerObjectiveHandlers(s *crud.Store, mux *http.ServeMux) {
	oh := objectiveHandler(s, s.RunTx)
	mux.HandleFunc("/o/", oh)
	mux.HandleFunc("/o", oh)
}

func objectiveHandler(s crud.ObjectiveStore, runTx crud.TxRun) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		klog.V(11).Info(r.Method, " ", r.URL.Path)

		switch {
		case r.Method == http.MethodPut:
			handleObjectivePUT(s, runTx, w, r)

		default:
			klog.Error("Not allowed ", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

type Objective struct {
	HItem
}

func createObjectiveEntity(ctx context.Context, s crud.ObjectiveStore, runTx crud.TxRun, o Objective) (crud.ObjectiveEntity, error) {
	data := crud.Objective{crud.HItem{crud.Item{o.Summary, o.Description}, o.ParentID}}
	
	if data.ParentID == "" {
		return s.CreateObjective(ctx, nil, data)
	}

	var entity crud.ObjectiveEntity
	err := runTx(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		_, err := s.ReadObjective(ctx, tx, data.ParentID)
		if err != nil {
			return err
		}

		entity, err = s.CreateObjective(ctx, tx, data)
		return err
	})

	return entity, err
}

func handleObjectivePUT(s crud.ObjectiveStore, runTx crud.TxRun, w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "/o") {
		klog.Error("invalid path:", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var o Objective
	err := decodeRequestBody(r, &o)
	if err != nil {
		klog.Error("decode failed:", r)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if o.Summary == "" {
		klog.Errorf("summary not set:%+v", o)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	entity, err := createObjectiveEntity(ctx, s, runTx, o)
	if errs.IsNotFound(err) {
		klog.Errorf("Parent not found:%+v", o)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		klog.Errorf("createObjectiveEntity failed; err:%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	klog.Info("objective created:", entity.ID)
}
