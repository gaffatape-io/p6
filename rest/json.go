package rest

import (
	"encoding/json"
	"k8s.io/klog"
	"net/http"
)

func decodeRequestBody(r *http.Request, data interface{}) error {
	klog.V(11).Infof("decodeRequest:%+v", r)
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&data)
	return err
}
