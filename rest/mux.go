package rest

import (
	"github.com/gaffatape-io/p6/crud"
	"github.com/gaffatape-io/p6/okrs"
	"net/http"
)

// API behavior:
// single objective and key results are mapped into REST resources.
// /o/<objective-id> => returns a single objective
// /k/<key-result-id> => returns a single key result
//
// Both objectives and key results may have parents and childs.
// If they do then they contain references to their respective.
//
// TODO: teams, organization and views

func keyResultHandler(s *crud.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("kr")
	}
}

// NewMux creates a new ServeMux for the rest API.
func NewMux(s *crud.Store, o *okrs.Objectives) *http.ServeMux {
	mux := http.NewServeMux()
	registerObjectiveHandlers(o, mux)
	mux.HandleFunc("/k/", keyResultHandler(s))
	return mux
}
