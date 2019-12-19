package rest

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"path"
)

func registerObjectiveHandlers(c *firestore.Client, mux *http.ServeMux) {
	oh := objectiveHandler(c)
	mux.HandleFunc("/o/", oh)
	mux.HandleFunc("/o", oh)
}

func objectiveHandler(c *firestore.Client) http.HandlerFunc {
	objectives := c.Collection("objectives")

	return func(w http.ResponseWriter, r *http.Request) {

		switch {
		case r.Method == http.MethodPut:
			handleObjectivePUT(objectives, w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func handleObjectivePUT(objectives *firestore.CollectionRef, w http.ResponseWriter, r *http.Request) {
	var o Objective
	err := decodeRequest(r, &o)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id := path.Base(r.URL.Path)
	if id != "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, _, err = objectives.Add(r.Context(), o)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
