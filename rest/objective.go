package rest

import (
	"context"
	"cloud.google.com/go/firestore"
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

func createObjective(ctx context.Context, s crud.ObjectiveStore, tx *firestore.Transaction, o Objective) (Objective, error) {
	ob := crud.Objective{crud.HItem{crud.Item{o.Summary, o.Description}, o.ParentID}}
	obe, err := s.CreateObjective(ctx, tx, ob)
	return Objective{HItem{Item{obe.ID, obe.Summary, obe.Description}, obe.ParentID}}, err
}

func handleObjectivePUT(s crud.ObjectiveStore, runTx crud.TxRun, w http.ResponseWriter, r *http.Request) {
	var objNew Objective
	err := decodeRequestBody(r, &objNew)
	if err != nil {
		klog.Error("decode failed:", r)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !strings.HasSuffix(r.URL.Path, "/o") {
		klog.Error("invalid path:", r.URL.Path)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if objNew.Summary == "" {
		klog.Errorf("summary not set:%+v", objNew)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	objCreated, err := createObjective(r.Context(), s, nil, objNew)
	if err != nil {
		klog.Errorf("createObjective failed; err:%+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	klog.Info("objective created:", objCreated.ID)
}
