package rest

import (
	"encoding/json"
	"net/http"
)


func decodeRequest(r *http.Request, data interface{}) error {
	dec := json.NewDecoder(r.Body)
	var o Objective
	return dec.Decode(&o)

}
