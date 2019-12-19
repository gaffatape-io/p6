package rest

import (
	"net/http"
	"testing"
)

func TestObjectiveAllowedMethods(t *testing.T) {
	checkMethodAllowed(t, "/o/<ignored>", http.MethodPut)
}

func TestObjectiveHandlerPUT(t *testing.T) {
	o := &Objective{/*Name: "a", Desc: "b"*/}
	resp := checkOK(t, roundTripPUT(t, "/o", o))
	t.Log(resp)
}
