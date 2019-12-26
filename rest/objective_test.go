package rest

import (
	"context"
	"net/http"
	"reflect"
	"testing"
)

func TestObjectiveAllowedMethods(t *testing.T) {
	RunRestTest(t, func(ctx context.Context, e *RestEnv) {
		e.checkMethodAllowed(ctx, "/o/<ignored>", http.MethodPut)
	})
}

func TestObjectiveHandlerPUT(t *testing.T) {
	RunRestTest(t, func(ctx context.Context, e *RestEnv) {
		sum := e.String("SsSs")
		desc := e.String("DdDd")
		o := &Objective{HItem{Item{Summary: sum, Description: desc}, ""}}
		resp := checkOK(t, e.roundTripPUT(ctx, "/o", o))
		t.Log(resp)

		objs := e.Firestore.Collection("objectives")
		matches, err := objs.Where("Summary", "==", sum).Documents(ctx).GetAll()
		t.Log(matches, err)
		if err != nil || len(matches) != 1 {
			t.Fatal()
		}

		want := map[string]interface{}{
			"Summary":     sum,
			"Description": desc,
			"ParentID":    "",
		}

		got := matches[0].Data()
		t.Log(got, want)
		if !reflect.DeepEqual(matches[0].Data(), want) {
			t.Fatal()
		}
	})
}

func TestObjectiveHandlerPUT_summaryMissing(t *testing.T) {
	RunRestTest(t, func(ctx context.Context, e *RestEnv) {
		desc := e.String("DDDD")
		o := &Objective{HItem{Item{Description: desc}, ""}}
		resp := checkResponse(t, http.StatusBadRequest, e.roundTripPUT(ctx, "/o", o))
		t.Log(resp)

		objs := e.Firestore.Collection("objectives")
		t.Log(objs)

	})
}
