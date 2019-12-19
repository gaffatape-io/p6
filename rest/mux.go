package rest

import (
	"net/http"

	"cloud.google.com/go/firestore"
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


func keyResultHandler(c *firestore.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic("kr")
	}
}

// NewMux creates a new ServeMux for the rest API.
func NewMux(c *firestore.Client) *http.ServeMux {	
	mux := http.NewServeMux()
	registerObjectiveHandlers(c, mux)
	mux.HandleFunc("/k/", keyResultHandler(c))
	return mux
}
